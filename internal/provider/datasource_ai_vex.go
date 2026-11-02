package provider

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
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

// vulnetix_vex_statements reads the tenant's published OpenVEX statements
// from the review decisions ledger (GET /review/decisions?kind=vex).

type vexStatementsDataSource struct{ client *client.Client }

func newVexStatementsDataSource() datasource.DataSource { return &vexStatementsDataSource{} }

type vexStatementsModel struct {
	Status     types.String        `tfsdk:"status"`
	Query      types.String        `tfsdk:"query"`
	Limit      types.Int64         `tfsdk:"limit"`
	Statements []vexStatementModel `tfsdk:"statements"`
}

type vexStatementModel struct {
	UUID              types.String `tfsdk:"uuid"`
	VulnerabilityID   types.String `tfsdk:"vulnerability_id"`
	VulnerabilityName types.String `tfsdk:"vulnerability_name"`
	Status            types.String `tfsdk:"status"`
	Justification     types.String `tfsdk:"justification"`
	Detail            types.String `tfsdk:"detail"`
	Author            types.String `tfsdk:"author"`
	Timestamp         types.Int64  `tfsdk:"timestamp"`
}

// defaultVexLimit bounds a read that sets no limit: a tenant can hold far
// more statements than a plan should carry.
const defaultVexLimit = 1000

func (d *vexStatementsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vex_statements"
}

func (d *vexStatementsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The tenant's published OpenVEX statements, newest first, as the review decisions " +
			"ledger lists them (`GET /review/decisions?kind=vex`). Statements still `under_investigation` are " +
			"not listed. Statements come from triage and from `vulnetix_vex_source` imports.",
		Attributes: map[string]schema.Attribute{
			"status": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Only statements with this status: `not_affected`, `affected` or `fixed`.",
				Validators:          []validator.String{stringvalidator.OneOf("not_affected", "affected", "fixed")},
			},
			"query": schema.StringAttribute{Optional: true, MarkdownDescription: "Only statements whose vulnerability, justification or detail contains this, ignoring case."},
			"limit": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "The most statements to read. Defaults to 1000.",
				Validators:          []validator.Int64{int64validator.Between(1, 100000)},
			},
			"statements": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The matching statements, newest first.",
				NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
					"uuid":               schema.StringAttribute{Computed: true, MarkdownDescription: "The statement's id."},
					"vulnerability_id":   schema.StringAttribute{Computed: true, MarkdownDescription: "The vulnerability, such as `CVE-2024-3094`."},
					"vulnerability_name": schema.StringAttribute{Computed: true, MarkdownDescription: "The vulnerability's name, when the statement gives one."},
					"status":             schema.StringAttribute{Computed: true, MarkdownDescription: "`not_affected`, `affected` or `fixed`."},
					"justification":      schema.StringAttribute{Computed: true, MarkdownDescription: "The VEX justification for `not_affected`, such as `vulnerable_code_not_in_execute_path`."},
					"detail":             schema.StringAttribute{Computed: true, MarkdownDescription: "The status notes, or the impact statement."},
					"author":             schema.StringAttribute{Computed: true, MarkdownDescription: "The VEX document's author."},
					"timestamp":          schema.Int64Attribute{Computed: true, MarkdownDescription: "The statement's timestamp, in Unix milliseconds."},
				}},
			},
		},
	}
}

func (d *vexStatementsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = dataSourceClient(req, resp)
}

func (d *vexStatementsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError("Provider not configured", "The vulnetix provider was not configured.")
		return
	}
	var m vexStatementsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &m)...)
	if resp.Diagnostics.HasError() {
		return
	}
	limit := defaultVexLimit
	if !m.Limit.IsNull() && !m.Limit.IsUnknown() {
		limit = int(m.Limit.ValueInt64())
	}
	sts, err := d.client.ListVexStatements(ctx, client.VexStatementFilter{
		Status: m.Status.ValueString(), Query: m.Query.ValueString(), Limit: limit,
	})
	if err != nil {
		resp.Diagnostics.AddError("Could not list VEX statements", err.Error())
		return
	}
	m.Statements = []vexStatementModel{}
	for _, s := range sts {
		id, name, _ := strings.Cut(s.Label, " · ")
		m.Statements = append(m.Statements, vexStatementModel{
			UUID:              types.StringValue(s.UUID),
			VulnerabilityID:   types.StringValue(id),
			VulnerabilityName: types.StringValue(name),
			Status:            types.StringValue(s.Status),
			Justification:     types.StringValue(s.Justification),
			Detail:            types.StringValue(s.Detail),
			Author:            types.StringValue(s.Author),
			Timestamp:         types.Int64Value(s.At),
		})
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &m)...)
}
