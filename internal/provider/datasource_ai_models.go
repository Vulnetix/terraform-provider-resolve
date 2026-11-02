package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/vulnetix/terraform-provider-vulnetix/internal/client"
)

// vulnetix_ai_models reads the AI firewall model catalogue
// (GET /ai-firewall/models): the built-in providers' models and the
// tenant's custom providers'. Use it to pin models in a
// vulnetix_ai_firewall_policy by slug.

type aiModelsDataSource struct{ client *client.Client }

func newAiModelsDataSource() datasource.DataSource { return &aiModelsDataSource{} }

type aiModelsModel struct {
	ProviderUUID types.String   `tfsdk:"provider_uuid"`
	Search       types.String   `tfsdk:"search"`
	Models       []aiModelModel `tfsdk:"models"`
}

type aiModelModel struct {
	UUID         types.String `tfsdk:"uuid"`
	Slug         types.String `tfsdk:"slug"`
	Name         types.String `tfsdk:"name"`
	ProviderUUID types.String `tfsdk:"provider_uuid"`
	ProviderSlug types.String `tfsdk:"provider_slug"`
	IsActive     types.Bool   `tfsdk:"is_active"`
	ReleasedAt   types.Int64  `tfsdk:"released_at"`
}

func (d *aiModelsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ai_models"
}

func (d *aiModelsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The AI models the AI firewall knows (`GET /ai-firewall/models`): every built-in provider's " +
			"models and those of the tenant's custom providers (`vulnetix_ai_provider`). Use a model's " +
			"`provider_slug` and `slug` to pin it in `vulnetix_ai_firewall_policy`.",
		Attributes: map[string]schema.Attribute{
			"provider_uuid": schema.StringAttribute{Optional: true, MarkdownDescription: "Only this provider's models, by the provider's uuid."},
			"search":        schema.StringAttribute{Optional: true, MarkdownDescription: "Only models whose slug or name contains this, ignoring case."},
			"models": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The matching models.",
				NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
					"uuid":          schema.StringAttribute{Computed: true, MarkdownDescription: "The model's id."},
					"slug":          schema.StringAttribute{Computed: true, MarkdownDescription: "The model's slug, as requests name it."},
					"name":          schema.StringAttribute{Computed: true, MarkdownDescription: "The model's display name."},
					"provider_uuid": schema.StringAttribute{Computed: true, MarkdownDescription: "The provider's id."},
					"provider_slug": schema.StringAttribute{Computed: true, MarkdownDescription: "The provider's slug, such as `anthropic` or a custom provider's name."},
					"is_active":     schema.BoolAttribute{Computed: true, MarkdownDescription: "Whether the provider still serves the model."},
					"released_at":   schema.Int64Attribute{Computed: true, MarkdownDescription: "When the model was released, in Unix milliseconds; 0 when unknown."},
				}},
			},
		},
	}
}

func (d *aiModelsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = dataSourceClient(req, resp)
}

func (d *aiModelsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError("Provider not configured", "The vulnetix provider was not configured.")
		return
	}
	var m aiModelsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &m)...)
	if resp.Diagnostics.HasError() {
		return
	}
	models, err := d.client.ListAiModels(ctx, client.AiModelFilter{
		ProviderUUID: m.ProviderUUID.ValueString(),
		Search:       m.Search.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Could not list AI models", err.Error())
		return
	}
	m.Models = []aiModelModel{}
	for _, x := range models {
		m.Models = append(m.Models, aiModelModel{
			UUID:         types.StringValue(x.UUID),
			Slug:         types.StringValue(x.Slug),
			Name:         types.StringValue(x.Name),
			ProviderUUID: types.StringValue(x.ProviderUUID),
			ProviderSlug: types.StringValue(x.ProviderSlug),
			IsActive:     types.BoolValue(x.IsActive),
			ReleasedAt:   types.Int64Value(x.ReleasedAt),
		})
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &m)...)
}
