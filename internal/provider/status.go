package provider

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"sort"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/vulnetix/terraform-provider-resolve/internal/client"
)

// Every resource has a computed status block: what the server says about the
// object beyond its spec. Most kinds are simply Applied. Some are applied and
// then wait:
//
//   - Pending: a Member invited and not yet accepted;
//   - AwaitingAuthorization: a messaging connection, such as Slack, that a
//     person has to authorize in a browser at status.authorization_url.
//
// Both are success: the object exists, and Terraform does not wait for the
// person. A refresh reads the phase again, so a plan shows when it changes.
// Failed is an error, not a state: a SIEM kind is written in Guardian after
// the config database, and when Guardian refuses, nothing is recorded and
// the apply fails with the reason.
//
// Some kinds show a value once, when the object is created: a ServiceAccount's
// API key and secret, a messaging connection's signing secret, a SIEM log
// source's ingest token. Terraform has no write-only computed attribute, so
// they are in status.once, sensitive, and kept in state from the create (the
// server never shows them again, and a refresh keeps what the state holds).
// Keep the state encrypted.

var (
	statusExtraType = tftypes.Map{ElementType: tftypes.String}
	statusAttrTypes = map[string]tftypes.Type{
		"phase":             tftypes.String,
		"message":           tftypes.String,
		"authorization_url": tftypes.String,
		"object_uuid":       tftypes.String,
		"extra":             statusExtraType,
		"once":              statusExtraType,
	}
	statusType = tftypes.Object{AttributeTypes: statusAttrTypes}

	// secretKey matches the status fields that are credentials: they go to
	// status.once, sensitive.
	secretKey = regexp.MustCompile(`(?i)secret|token|key|password|credential`)
)

func statusAttribute() schema.Attribute {
	return schema.SingleNestedAttribute{
		MarkdownDescription: "What the server says about the object beyond its spec.",
		Computed:            true,
		Attributes: map[string]schema.Attribute{
			"phase": schema.StringAttribute{
				Computed: true,
				MarkdownDescription: "`Applied` for most objects. `Pending` when the object waits on someone (an invited member who has not accepted), " +
					"`AwaitingAuthorization` when a person has to authorize it in a browser at `authorization_url`. Terraform does not wait for either. " +
					"`Unmanaged` means the object exists but no manager owns it.",
			},
			"message":           schema.StringAttribute{Computed: true, MarkdownDescription: "Why the object is in its phase, when the server says."},
			"authorization_url": schema.StringAttribute{Computed: true, MarkdownDescription: "Where a person authorizes the object while the phase is `AwaitingAuthorization`."},
			"object_uuid":       schema.StringAttribute{Computed: true, MarkdownDescription: "The object's id on the server."},
			"extra": schema.MapAttribute{
				Computed: true, ElementType: types.StringType,
				MarkdownDescription: "Other fields the kind reports, such as an invitation's expiry or a service account's principal. Structured values are JSON.",
			},
			"once": schema.MapAttribute{
				Computed: true, ElementType: types.StringType, Sensitive: true,
				MarkdownDescription: "Credentials the server shows only when it creates the object: a service account's `apiKey` and `secret`, a messaging " +
					"connection's `signingSecret`, a log source's ingest token. They are kept in the state from the create, so keep the state encrypted. " +
					"They cannot be read again: replace the object to get new ones.",
			},
		},
	}
}

// statusScalar is a status value as a string. Structured values are JSON.
func statusScalar(v any) string {
	switch x := v.(type) {
	case string:
		return x
	case nil:
		return ""
	case fmt.Stringer:
		return x.String()
	case bool, int, int64, float64:
		return fmt.Sprint(x)
	}
	return jsonString(v)
}

// buildStatus is the status block for what the server answered. once is the
// block's previous once map, kept when the answer has none (an update, a
// refresh): the server shows those values once.
func buildStatus(st client.ObjectStatus, prior tftypes.Value) tftypes.Value {
	str := func(s string) tftypes.Value {
		if s == "" {
			return tftypes.NewValue(tftypes.String, nil)
		}
		return tftypes.NewValue(tftypes.String, s)
	}
	extra, once := map[string]tftypes.Value{}, map[string]tftypes.Value{}
	authURL := ""
	keys := make([]string, 0, len(st.Extra))
	for k := range st.Extra {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		v := statusScalar(st.Extra[k])
		switch {
		case k == "authorizationUrl":
			authURL = v
		case secretKey.MatchString(k):
			once[k] = tftypes.NewValue(tftypes.String, v)
		default:
			extra[k] = tftypes.NewValue(tftypes.String, v)
		}
	}
	vals := map[string]tftypes.Value{
		"phase":             str(st.Phase),
		"message":           str(st.Message),
		"authorization_url": str(authURL),
		"object_uuid":       str(st.ObjectUUID),
		"extra":             tftypes.NewValue(statusExtraType, nil),
		"once":              tftypes.NewValue(statusExtraType, nil),
	}
	if len(extra) > 0 {
		vals["extra"] = tftypes.NewValue(statusExtraType, extra)
	}
	switch {
	case len(once) > 0:
		vals["once"] = tftypes.NewValue(statusExtraType, once)
	case isSet(prior):
		if pm, err := objectAttrs(prior); err == nil && isSet(nullOr(pm["once"])) {
			vals["once"] = pm["once"]
		}
	}
	// A refresh does not carry what only an apply reports (the authorization
	// URL, an invitation's expiry): while the phase is the one the apply
	// left, keep what it said.
	if pm, err := objectAttrs(prior); err == nil && isSet(nullOr(pm["phase"])) {
		var pp string
		if pm["phase"].As(&pp) == nil && pp == st.Phase {
			if authURL == "" {
				vals["authorization_url"] = pm["authorization_url"]
			}
			if len(extra) == 0 && isSet(nullOr(pm["extra"])) {
				vals["extra"] = pm["extra"]
			}
		}
	}
	return tftypes.NewValue(statusType, vals)
}

// failedStatus returns the error a Failed status stands for, or "".
func failedStatus(st client.ObjectStatus) string {
	if st.Phase != client.PhaseFailed {
		return ""
	}
	if st.Message != "" {
		return st.Message
	}
	return "the server reported the write as Failed"
}

// ownerHint adds what to do about a refusal that is about the caller's role,
// for the kinds only an owner may apply.
func ownerHint(spec kindSpec, err error) string {
	var ae *client.APIError
	if !errors.As(err, &ae) || !spec.OwnerOnly {
		return ""
	}
	if ae.Status == 403 || ae.Code == "forbidden" || ae.Code == "kind_not_allowed" {
		return "\n\n" + spec.Kind + " decides who can act for the tenant, so it can only be applied by an owner or admin of the organisation, " +
			"and a role above the caller's own cannot be granted. Authenticate with the ApiKey of an admin or owner service account, " +
			"or with a bearer token for an owner or admin."
	}
	return ""
}

// guardianHint explains a 502 from an external kind.
func guardianHint(err error) string {
	var ae *client.APIError
	if errors.As(err, &ae) && ae.Code == "external_failed" {
		return "\n\nThis kind is written in Guardian after the config database. Guardian refused or failed the write, and nothing was recorded; " +
			"fix the reason above and apply again."
	}
	return ""
}

func jsonString(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprint(v)
	}
	return string(b)
}
