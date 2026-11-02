package provider

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/vulnetix/terraform-provider-vulnetix/internal/client"
)

// Attributes every kind resource has besides its spec.
const (
	attrID          = "id"
	attrName        = "name"
	attrLabels      = "labels"
	attrAnnotations = "annotations"
	attrStatus      = "status"
)

// reservedAttrs may not be generated from a spec: the first four are the
// resource's own, the rest are Terraform meta-arguments.
var reservedAttrs = map[string]bool{
	attrID: true, attrName: true, attrLabels: true, attrAnnotations: true, attrStatus: true,
	"count": true, "for_each": true, "provider": true, "lifecycle": true,
	"depends_on": true, "provisioner": true, "connection": true,
}

// commonNamePattern is common.schema.json's name pattern. A kind may
// narrow it (AiProvider's slug), and then the message quotes the pattern.
const commonNamePattern = "^[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?$"

// privateImported marks state that came from ImportState, so the next Read
// fills in every field the server has rather than only the managed ones.
const privateImported = "imported"

var (
	_ resource.ResourceWithConfigure   = (*kindResource)(nil)
	_ resource.ResourceWithImportState = (*kindResource)(nil)
)

// kindResource is the resource for one config kind. Every kind uses this
// type; only the kindSpec differs.
type kindResource struct {
	spec   kindSpec
	client *client.Client
}

func newKindResource(spec kindSpec) func() resource.Resource {
	return func() resource.Resource { return &kindResource{spec: spec} }
}

func (r *kindResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + r.spec.TypeName
}

func (r *kindResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	attrs := r.spec.Attributes()
	attrs[attrID] = schema.StringAttribute{
		MarkdownDescription: "The object's name, the same as `name`.",
		Computed:            true,
		PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
	}
	if r.spec.Singleton {
		attrs[attrName] = schema.StringAttribute{
			MarkdownDescription: "The object's name. " + r.spec.Kind + " is one per tenant, so it is always `default`.",
			Optional:            true,
			Computed:            true,
			Default:             stringdefault.StaticString("default"),
			Validators:          []validator.String{stringvalidator.OneOf("default")},
			PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
		}
	} else {
		v := []validator.String{stringvalidator.LengthAtMost(63)}
		if r.spec.NamePattern != "" {
			msg := "must match " + r.spec.NamePattern
			if r.spec.NamePattern == commonNamePattern {
				msg = "must be lower case letters, digits and hyphens, starting and ending with a letter or digit"
			}
			v = append(v, stringvalidator.RegexMatches(regexp.MustCompile(r.spec.NamePattern), msg))
		}
		attrs[attrName] = schema.StringAttribute{
			MarkdownDescription: "The object's name (`metadata.name`), unique per kind within the tenant. Changing it replaces the object.",
			Required:            true,
			Validators:          v,
			PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
		}
	}
	attrs[attrStatus] = statusAttribute()
	attrs[attrLabels] = schema.MapAttribute{
		MarkdownDescription: "Free key/value labels (`metadata.labels`). Stored with the object, never interpreted by the server.",
		ElementType:         types.StringType,
		Optional:            true,
	}
	attrs[attrAnnotations] = schema.MapAttribute{
		MarkdownDescription: "Free key/value notes for people and tools (`metadata.annotations`), such as a ticket link.",
		ElementType:         types.StringType,
		Optional:            true,
	}
	deps := ""
	if len(r.spec.DependsOn) > 0 {
		names := make([]string, len(r.spec.DependsOn))
		for i, k := range r.spec.DependsOn {
			names[i] = "`" + k + "`"
		}
		deps = "\n\nAn object of this kind can reference " + strings.Join(names, ", ") +
			" objects. Reference those resources' attributes, or add `depends_on`, so Terraform applies them first."
	}
	if r.spec.OwnerOnly {
		deps += "\n\n" + r.spec.Kind + " decides who can act for the tenant, so it can only be applied with a token that has the owner role and lists `" + r.spec.Kind + "` among its kinds."
	}
	resp.Schema = schema.Schema{
		MarkdownDescription: r.spec.Description + deps + "\n\nManages the `" + r.spec.Kind + "` config kind through `/config/objects/" + r.spec.Kind + "/{name}`. " +
			"Applies are recorded with the provider's `manager`, so a GitOps `vulnetix config apply --prune` leaves this object alone.\n\n" +
			"Most kinds merge: an argument removed from the configuration is no longer managed, and keeps its last value on the server. A kind that is wholly declared (most of them, such as `ApiSpecSource`, `RiskStrategy` and `VexSource`) is replaced instead: a removed argument goes back to its default. Destroying the resource deletes the object, or lets go of it for a kind the server never deletes (such as `Organization`). Some kinds wait on a person after they are applied (an invited `Member`, a Slack connection): `status.phase` says so, and Terraform does not wait.",
		Attributes: attrs,
	}
}

func (r *kindResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected provider data", fmt.Sprintf("expected *client.Client, got %T", req.ProviderData))
		return
	}
	r.client = c
}

// document builds the document to PUT from a plan.
func (r *kindResource) document(plan, config, prior tftypes.Value) (client.Document, error) {
	attrs, err := objectAttrs(plan)
	if err != nil {
		return client.Document{}, err
	}
	var name string
	if err := attrs[attrName].As(&name); err != nil {
		return client.Document{}, fmt.Errorf("name: %w", err)
	}
	labels, err := stringMap(attrs[attrLabels])
	if err != nil {
		return client.Document{}, fmt.Errorf("labels: %w", err)
	}
	annotations, err := stringMap(attrs[attrAnnotations])
	if err != nil {
		return client.Document{}, fmt.Errorf("annotations: %w", err)
	}
	spec, err := toSpec(specInput{plan: plan, config: config, prior: prior, create: !isSet(prior)}, r.spec.Fields)
	if err != nil {
		return client.Document{}, err
	}
	return client.Document{
		APIVersion: client.APIVersion,
		Kind:       r.spec.Kind,
		Metadata:   client.Metadata{Name: name, Labels: labels, Annotations: annotations},
		Spec:       spec,
	}, nil
}

func stringMap(v tftypes.Value) (map[string]string, error) {
	if !isSet(nullOr(v)) {
		return nil, nil
	}
	var m map[string]tftypes.Value
	if err := v.As(&m); err != nil {
		return nil, err
	}
	out := make(map[string]string, len(m))
	for k, e := range m {
		var s string
		if err := e.As(&s); err != nil {
			return nil, err
		}
		out[k] = s
	}
	return out, nil
}

// stateFromPlan is the state after a successful apply: the plan with id
// filled in. The plan, not the server's answer, because Terraform requires
// the new state to equal every configured value, and the server's answer can
// carry fields the configuration does not set. Write-only values are already
// null in a plan.
func stateFromPlan(plan tftypes.Value, status tftypes.Value) (tftypes.Value, error) {
	attrs, err := objectAttrs(plan)
	if err != nil {
		return tftypes.Value{}, err
	}
	attrs[attrID] = tftypes.NewValue(tftypes.String, nil)
	if n, ok := attrs[attrName]; ok {
		attrs[attrID] = n
	}
	attrs[attrStatus] = status
	return tftypes.NewValue(plan.Type(), attrs), nil
}

func (r *kindResource) apply(ctx context.Context, plan, config, prior tftypes.Value, diags *diag.Diagnostics) (tftypes.Value, bool) {
	if r.client == nil {
		diags.AddError("Provider not configured", "The vulnetix provider was not configured.")
		return tftypes.Value{}, false
	}
	doc, err := r.document(plan, config, prior)
	if err != nil {
		diags.AddError("Could not build the "+r.spec.Kind+" document", err.Error())
		return tftypes.Value{}, false
	}
	obj, err := r.client.PutObject(ctx, doc)
	if err != nil {
		diags.AddError("Could not apply "+r.spec.Kind+"/"+doc.Metadata.Name, err.Error()+ownerHint(r.spec, err)+guardianHint(err))
		return tftypes.Value{}, false
	}
	if msg := failedStatus(obj.Status); msg != "" {
		diags.AddError("Could not apply "+r.spec.Kind+"/"+doc.Metadata.Name, msg)
		return tftypes.Value{}, false
	}
	var priorStatus tftypes.Value
	if pa, err := objectAttrs(prior); err == nil {
		priorStatus = nullOr(pa[attrStatus])
	}
	state, err := stateFromPlan(plan, buildStatus(obj.Status, priorStatus))
	if err != nil {
		diags.AddError("Could not record state", err.Error())
		return tftypes.Value{}, false
	}
	return state, true
}

func (r *kindResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	state, ok := r.apply(ctx, req.Plan.Raw, req.Config.Raw, nullValue, &resp.Diagnostics)
	if ok {
		resp.State.Raw = state
	}
}

func (r *kindResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	state, ok := r.apply(ctx, req.Plan.Raw, req.Config.Raw, req.State.Raw, &resp.Diagnostics)
	if ok {
		resp.State.Raw = state
	}
}

func (r *kindResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Provider not configured", "The vulnetix provider was not configured.")
		return
	}
	prior := req.State.Raw
	attrs, err := objectAttrs(prior)
	if err != nil {
		resp.Diagnostics.AddError("Could not read state", err.Error())
		return
	}
	var name string
	if err := attrs[attrName].As(&name); err != nil || name == "" {
		resp.Diagnostics.AddError("State has no name", "The "+r.spec.Kind+" state has no name to read.")
		return
	}
	imported, d := req.Private.GetKey(ctx, privateImported)
	resp.Diagnostics.Append(d...)
	keepAll := len(imported) > 0

	obj, err := r.client.GetObject(ctx, r.spec.Kind, name)
	if errors.Is(err, client.ErrNotFound) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Could not read "+r.spec.Kind+"/"+name, err.Error())
		return
	}
	typ, ok := prior.Type().(tftypes.Object)
	if !ok {
		resp.Diagnostics.AddError("Unexpected state type", prior.Type().String())
		return
	}
	vals, err := fromSpec(obj.Document.Spec, typ, prior, r.spec.Fields, keepAll)
	if err != nil {
		resp.Diagnostics.AddError("Could not read "+r.spec.Kind+"/"+name, err.Error())
		return
	}
	vals[attrID] = tftypes.NewValue(tftypes.String, name)
	vals[attrName] = tftypes.NewValue(tftypes.String, name)
	vals[attrLabels] = metaMap(obj.Document.Metadata.Labels, attrs[attrLabels], keepAll)
	vals[attrAnnotations] = metaMap(obj.Document.Metadata.Annotations, attrs[attrAnnotations], keepAll)
	vals[attrStatus] = buildStatus(obj.Status, nullOr(attrs[attrStatus]))
	resp.State.Raw = tftypes.NewValue(typ, vals)
	if keepAll {
		resp.Diagnostics.Append(resp.Private.SetKey(ctx, privateImported, nil)...)
	}
}

// metaMap is labels or annotations as state. An empty map from the server
// reads as null unless the configuration wrote an empty map.
func metaMap(m map[string]string, prior tftypes.Value, keepAll bool) tftypes.Value {
	mt := tftypes.Map{ElementType: tftypes.String}
	if len(m) == 0 {
		if isSet(nullOr(prior)) {
			return tftypes.NewValue(mt, map[string]tftypes.Value{})
		}
		return tftypes.NewValue(mt, nil)
	}
	if !keepAll && !isSet(nullOr(prior)) {
		return tftypes.NewValue(mt, nil)
	}
	vals := make(map[string]tftypes.Value, len(m))
	for k, v := range m {
		vals[k] = tftypes.NewValue(tftypes.String, v)
	}
	return tftypes.NewValue(mt, vals)
}

func (r *kindResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Provider not configured", "The vulnetix provider was not configured.")
		return
	}
	var name string
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root(attrName), &name)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.DeleteObject(ctx, r.spec.Kind, name); err != nil {
		resp.Diagnostics.AddError("Could not delete "+r.spec.Kind+"/"+name, err.Error()+ownerHint(r.spec, err)+guardianHint(err))
	}
}

// ImportState imports by name: terraform import vulnetix_quality_gate.x default.
func (r *kindResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if req.ID == "" {
		resp.Diagnostics.AddError("Missing name", "Import a "+r.spec.Kind+" by its name.")
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root(attrName), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root(attrID), req.ID)...)
	resp.Diagnostics.Append(resp.Private.SetKey(ctx, privateImported, []byte("true"))...)
}
