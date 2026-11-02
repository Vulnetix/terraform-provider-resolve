// Package provider is the vulnetix Terraform and OpenTofu provider.
package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/vulnetix/terraform-provider-vulnetix/internal/client"
)

// Environment variables the provider reads when an argument is not set.
const (
	EnvToken    = "VULNETIX_CONFIG_TOKEN"
	EnvEndpoint = "VULNETIX_ENDPOINT"
)

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
			"token": schema.StringAttribute{
				MarkdownDescription: "A config management token (`vmt_...`). Defaults to `$" + EnvToken + "`. " +
					"Applying needs a token with the owner or admin role.",
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
	token := firstSet(m.Token, os.Getenv(EnvToken))
	manager := firstSet(m.Manager, "")
	if token == "" {
		resp.Diagnostics.AddWarning("No Vulnetix token",
			"Set the provider's token argument or "+EnvToken+" to a config management token. Every request will be refused without one.")
	}
	c := client.New(client.Options{
		Endpoint:  endpoint,
		Token:     token,
		Manager:   manager,
		UserAgent: "terraform-provider-vulnetix/" + p.version,
	})
	resp.ResourceData = c
	resp.DataSourceData = c
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
