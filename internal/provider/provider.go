// Package provider is the vulnetix Terraform and OpenTofu provider.
package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/vulnetix/terraform-provider-resolve/internal/client"
)

// Environment variables the provider reads when an argument is not set.
const (
	EnvAPIKey   = "VULNETIX_API_KEY"
	EnvOrgID    = "VULNETIX_ORG_ID"
	EnvToken    = "VULNETIX_API_TOKEN"
	EnvEndpoint = "VULNETIX_ENDPOINT"
)

// roleNote is appended to every credential argument: what the principal must
// be, and what happens when it is not.
const roleNote = "The principal must be an owner or admin of the organisation: a service account " +
	"created with the admin or owner role is the recommended automation credential, and its ApiKey " +
	"the recommended way to authenticate. A member is refused with HTTP 403 on apply."

var _ provider.Provider = (*vulnetixProvider)(nil)

type vulnetixProvider struct {
	version string
	// kinds are the resources served. Tests add kinds generated from test
	// schemas; a release serves generatedKinds.
	kinds []kindSpec
}

// New returns the provider factory for a build version.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &vulnetixProvider{version: version, kinds: generatedKinds}
	}
}

type providerModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
	APIKey   types.String `tfsdk:"api_key"`
	OrgID    types.String `tfsdk:"org_id"`
	Token    types.String `tfsdk:"token"`
	Manager  types.String `tfsdk:"manager"`
}

func (p *vulnetixProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "vulnetix"
	resp.Version = p.version
}

func (p *vulnetixProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manage Vulnetix configuration as code. Every resource is one " +
			"`vulnetix.com/v1` config document, generated from the kind schemas in " +
			"[config-schemas](https://github.com/Vulnetix/config-schemas), and applied through the same " +
			"server-side engine as `vulnetix config apply`.",
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "The Vulnetix API base URL. Defaults to `$" + EnvEndpoint + "`, then `" + client.DefaultEndpoint + "`.",
				Optional:            true,
			},
			"api_key": schema.StringAttribute{
				MarkdownDescription: "The organisation ApiKey, sent as `Authorization: ApiKey <org_id>:<api_key>`. " +
					"Use it with `org_id`. Defaults to `$" + EnvAPIKey + "`. Set either `api_key` and `org_id`, or `token`, not both. " + roleNote,
				Optional:  true,
				Sensitive: true,
			},
			"org_id": schema.StringAttribute{
				MarkdownDescription: "The organisation UUID the `api_key` belongs to. Defaults to `$" + EnvOrgID + "`. " +
					"Required with `api_key`; not needed with `token`, which already names its organisation.",
				Optional: true,
			},
			"token": schema.StringAttribute{
				MarkdownDescription: "A bearer token, sent as `Authorization: Bearer <token>`: an Authentik OAuth2 " +
					"access token (for example from a client-credentials grant) or a Vulnetix session JWT. " +
					"Defaults to `$" + EnvToken + "`. Set either `token`, or `api_key` and `org_id`, not both. " + roleNote,
				Optional:  true,
				Sensitive: true,
			},
			"manager": schema.StringAttribute{
				MarkdownDescription: "The manager recorded as the owner of every object this provider applies. " +
					"A GitOps apply with `--prune` only deletes objects it manages itself. Defaults to `" + client.DefaultManager + "`.",
				Optional: true,
			},
		},
	}
}

func (p *vulnetixProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var m providerModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &m)...)
	if resp.Diagnostics.HasError() {
		return
	}
	endpoint := firstSet(m.Endpoint, os.Getenv(EnvEndpoint))
	manager := firstSet(m.Manager, "")
	apiKey := firstSet(m.APIKey, os.Getenv(EnvAPIKey))
	orgID := firstSet(m.OrgID, os.Getenv(EnvOrgID))
	token := firstSet(m.Token, os.Getenv(EnvToken))
	resp.Diagnostics.Append(checkCredentials(apiKey, orgID, token)...)
	if resp.Diagnostics.HasError() {
		return
	}
	c := client.New(client.Options{
		Endpoint:  endpoint,
		APIKey:    apiKey,
		OrgID:     orgID,
		Token:     token,
		Manager:   manager,
		UserAgent: "terraform-provider-resolve/" + p.version,
	})
	resp.ResourceData = c
	resp.DataSourceData = c
}

// checkCredentials validates the resolved credentials: either api_key with
// org_id, or token, never both, and a warning when there is none.
func checkCredentials(apiKey, orgID, token string) diag.Diagnostics {
	var d diag.Diagnostics
	switch {
	case apiKey != "" && token != "":
		d.AddError("Two Vulnetix credentials",
			"Both an ApiKey (api_key / "+EnvAPIKey+") and a bearer token (token / "+EnvToken+") are set. "+
				"Configure one: api_key with org_id, or token.")
	case apiKey != "" && orgID == "":
		d.AddError("api_key needs org_id",
			"An ApiKey is sent as \"ApiKey <org>:<key>\". Set org_id or "+EnvOrgID+" to the organisation the key belongs to.")
	case apiKey == "" && token == "" && orgID != "":
		d.AddError("org_id needs api_key",
			"org_id is set without an api_key. Set api_key or "+EnvAPIKey+", or use token or "+EnvToken+" instead.")
	case apiKey == "" && token == "":
		d.AddWarning("No Vulnetix credentials",
			"Set api_key and org_id ("+EnvAPIKey+" and "+EnvOrgID+"), or token ("+EnvToken+"). "+
				"Every request will be refused without one. The principal must be an owner or admin; "+
				"an admin or owner service account's ApiKey is the recommended automation credential.")
	}
	return d
}

// firstSet is the argument when it is set, else fallback. An unknown value
// (an argument computed from another resource) falls back too.
func firstSet(v types.String, fallback string) string {
	if !v.IsNull() && !v.IsUnknown() && v.ValueString() != "" {
		return v.ValueString()
	}
	return fallback
}

func (p *vulnetixProvider) Resources(context.Context) []func() resource.Resource {
	out := make([]func() resource.Resource, 0, len(p.kinds))
	for _, k := range p.kinds {
		out = append(out, newKindResource(k))
	}
	return out
}

func (p *vulnetixProvider) DataSources(context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		newKindsDataSource,
		newDetectionRuleFeedDataSource,
		newVirtualPatchesDataSource,
		newAiModelsDataSource,
		newVexStatementsDataSource,
		newGitHubInstallationReposDataSource,
	}
}
