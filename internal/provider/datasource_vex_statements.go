package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/vulnetix/terraform-provider-resolve/internal/client"
)

// vulnetix_vex_statements reads every VEX statement the tenant holds
// (GET /vex/statements), in every status, with its products and purls.

type vexStatementsDataSource struct{ client *client.Client }

func newVexStatementsDataSource() datasource.DataSource { return &vexStatementsDataSource{} }

type vexStatementsModel struct {
	Status          types.String        `tfsdk:"status"`
	VulnerabilityID types.String        `tfsdk:"vulnerability_id"`
	Purl            types.String        `tfsdk:"purl"`
	Repository      types.String        `tfsdk:"repository"`
	Limit           types.Int64         `tfsdk:"limit"`
	Statements      []vexStatementModel `tfsdk:"statements"`
}

type vexStatementModel struct {
	UUID              types.String `tfsdk:"uuid"`
	DocumentUUID      types.String `tfsdk:"document_uuid"`
	DocumentID        types.String `tfsdk:"document_id"`
	Author            types.String `tfsdk:"author"`
	VulnerabilityID   types.String `tfsdk:"vulnerability_id"`
	VulnerabilityName types.String `tfsdk:"vulnerability_name"`
	Aliases           []string     `tfsdk:"aliases"`
	Status            types.String `tfsdk:"status"`
	Justification     types.String `tfsdk:"justification"`
	StatusNotes       types.String `tfsdk:"status_notes"`
	ImpactStatement   types.String `tfsdk:"impact_statement"`
	ActionStatement   types.String `tfsdk:"action_statement"`
	Supplier          types.String `tfsdk:"supplier"`
	Products          types.String `tfsdk:"products"`
	Purls             []string     `tfsdk:"purls"`
	Timestamp         types.Int64  `tfsdk:"timestamp"`
	LastUpdated       types.Int64  `tfsdk:"last_updated"`
	Disposition       types.String `tfsdk:"disposition"`
	SuppressionUUID   types.String `tfsdk:"suppression_uuid"`
	Repository        types.String `tfsdk:"repository"`
	Grounded          types.Bool   `tfsdk:"grounded"`
	SourceFormat      types.String `tfsdk:"source_format"`
}

// defaultVexLimit bounds a read that sets no limit: a tenant can hold far
// more statements than a plan should carry.
const defaultVexLimit = 1000

var vexStatuses = []string{"not_affected", "affected", "fixed", "under_investigation"}

func (d *vexStatementsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vex_statements"
}

func (d *vexStatementsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	str := func(desc string) schema.StringAttribute {
		return schema.StringAttribute{Computed: true, MarkdownDescription: desc}
	}
	resp.Schema = schema.Schema{
		MarkdownDescription: "Every VEX statement the tenant holds, newest first (`GET /vex/statements`), across all its " +
			"documents and in every status, including `under_investigation` and `affected`. Statements come from triage, " +
			"from suppressions and from `vulnetix_vex_source` imports. A `vmt_` token needs the `VexSource` or " +
			"`Suppression` kind.",
		Attributes: map[string]schema.Attribute{
			"status": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Only statements with this status: `not_affected`, `affected`, `fixed` or `under_investigation`.",
				Validators:          []validator.String{stringvalidator.OneOf(vexStatuses...)},
			},
			"vulnerability_id": schema.StringAttribute{Optional: true, MarkdownDescription: "Only statements about this vulnerability, such as `CVE-2024-3094`, ignoring case."},
			"purl":             schema.StringAttribute{Optional: true, MarkdownDescription: "Only statements whose products contain this, such as `pkg:npm/lodash`."},
			"repository":       schema.StringAttribute{Optional: true, MarkdownDescription: "Only statements from this repository's documents, as `owner/name`."},
			"limit": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "The most statements to read. Defaults to 1000.",
				Validators:          []validator.Int64{int64validator.Between(1, 100000)},
			},
			"statements": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The matching statements, newest first.",
				NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
					"uuid":               str("The statement's id."),
					"document_uuid":      str("The VEX document's id."),
					"document_id":        str("The VEX document's own identifier (`@id`)."),
					"author":             str("The VEX document's author."),
					"vulnerability_id":   str("The vulnerability, such as `CVE-2024-3094`."),
					"vulnerability_name": str("The vulnerability's name, when the statement gives one."),
					"aliases":            schema.ListAttribute{Computed: true, ElementType: types.StringType, MarkdownDescription: "Other ids for the vulnerability."},
					"status":             str("`not_affected`, `affected`, `fixed` or `under_investigation`."),
					"justification":      str("The VEX justification for `not_affected`, such as `vulnerable_code_not_in_execute_path`."),
					"status_notes":       str("Notes on the status."),
					"impact_statement":   str("Why the product is not affected, in prose."),
					"action_statement":   str("What to do, for `affected`."),
					"supplier":           str("The supplier, when given."),
					"products":           str("The statement's products as stored (OpenVEX products or CycloneDX affects), as a JSON string."),
					"purls":              schema.ListAttribute{Computed: true, ElementType: types.StringType, MarkdownDescription: "The package URLs found in the products."},
					"timestamp":          schema.Int64Attribute{Computed: true, MarkdownDescription: "The statement's timestamp, in Unix milliseconds."},
					"last_updated":       schema.Int64Attribute{Computed: true, MarkdownDescription: "When the statement last changed, in Unix milliseconds; 0 when never."},
					"disposition":        str("How the statement was disposed, when it was."),
					"suppression_uuid":   str("The suppression the statement became, if any."),
					"repository":         str("The repository the document belongs to, as `owner/name`, if any."),
					"grounded":           schema.BoolAttribute{Computed: true, MarkdownDescription: "Whether the document is grounded in a scan of the product."},
					"source_format":      str("The document's format, such as `openvex` or `cyclonedx`."),
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
		Status:          m.Status.ValueString(),
		VulnerabilityID: m.VulnerabilityID.ValueString(),
		Purl:            m.Purl.ValueString(),
		Repository:      m.Repository.ValueString(),
		Limit:           limit,
	})
	if err != nil {
		resp.Diagnostics.AddError("Could not list VEX statements", err.Error())
		return
	}
	m.Statements = []vexStatementModel{}
	for _, s := range sts {
		products := "null"
		if len(s.Products) > 0 {
			products = string(s.Products)
		}
		m.Statements = append(m.Statements, vexStatementModel{
			UUID:              types.StringValue(s.UUID),
			DocumentUUID:      types.StringValue(s.DocumentUUID),
			DocumentID:        types.StringValue(s.DocumentID),
			Author:            types.StringValue(s.Author),
			VulnerabilityID:   types.StringValue(s.VulnerabilityID),
			VulnerabilityName: types.StringValue(s.VulnerabilityName),
			Aliases:           append([]string{}, s.Aliases...),
			Status:            types.StringValue(s.Status),
			Justification:     types.StringValue(s.Justification),
			StatusNotes:       types.StringValue(s.StatusNotes),
			ImpactStatement:   types.StringValue(s.ImpactStatement),
			ActionStatement:   types.StringValue(s.ActionStatement),
			Supplier:          types.StringValue(s.Supplier),
			Products:          types.StringValue(products),
			Purls:             append([]string{}, s.Purls...),
			Timestamp:         types.Int64Value(s.Timestamp),
			LastUpdated:       types.Int64Value(s.LastUpdated),
			Disposition:       types.StringValue(s.Disposition),
			SuppressionUUID:   types.StringValue(s.SuppressionUUID),
			Repository:        types.StringValue(s.RepositoryFullName),
			Grounded:          types.BoolValue(s.Grounded),
			SourceFormat:      types.StringValue(s.SourceFormat),
		})
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &m)...)
}
