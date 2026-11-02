package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/vulnetix/terraform-provider-vulnetix/internal/client"
)

type kindsDataSource struct{ client *client.Client }

func newKindsDataSource() datasource.DataSource { return &kindsDataSource{} }

type kindsModel struct {
	APIVersion    types.String `tfsdk:"api_version"`
	SchemaVersion types.String `tfsdk:"schema_version"`
	Kinds         []kindModel  `tfsdk:"kinds"`
	SchemaKinds   []string     `tfsdk:"schema_kinds"`
}

type kindModel struct {
	Kind      types.String `tfsdk:"kind"`
	SchemaURL types.String `tfsdk:"schema_url"`
	DependsOn []string     `tfsdk:"depends_on"`
	Schema    types.String `tfsdk:"schema"`
}

func (d *kindsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_kinds"
}

func (d *kindsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The config kinds the Vulnetix server applies, with their JSON Schemas (`GET /config/kinds`). " +
			"Use it to check that a kind is served before applying it.",
		Attributes: map[string]schema.Attribute{
			"api_version":    schema.StringAttribute{Computed: true, MarkdownDescription: "The document format version, `vulnetix.com/v1`."},
			"schema_version": schema.StringAttribute{Computed: true, MarkdownDescription: "The config-schemas release the server validates against."},
			"schema_kinds": schema.ListAttribute{
				Computed: true, ElementType: types.StringType,
				MarkdownDescription: "Every kind with a released schema, including kinds the server does not apply yet.",
			},
			"kinds": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The kinds the server applies, in dependency order.",
				NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
					"kind":       schema.StringAttribute{Computed: true, MarkdownDescription: "The kind, such as `QualityGate`."},
					"schema_url": schema.StringAttribute{Computed: true, MarkdownDescription: "Where the kind's released schema is downloaded from."},
					"depends_on": schema.ListAttribute{Computed: true, ElementType: types.StringType, MarkdownDescription: "Kinds applied before this one."},
					"schema":     schema.StringAttribute{Computed: true, MarkdownDescription: "The kind's JSON Schema, as a JSON string."},
				}},
			},
		},
	}
}

func (d *kindsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = dataSourceClient(req, resp)
}

func (d *kindsDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError("Provider not configured", "The vulnetix provider was not configured.")
		return
	}
	k, err := d.client.ListKinds(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Could not list kinds", err.Error())
		return
	}
	m := kindsModel{
		APIVersion:    types.StringValue(k.APIVersion),
		SchemaVersion: types.StringValue(k.SchemaVersion),
		Kinds:         []kindModel{},
		SchemaKinds:   append([]string{}, k.SchemaKinds...),
	}
	for _, kd := range k.Kinds {
		m.Kinds = append(m.Kinds, kindModel{
			Kind:      types.StringValue(kd.Kind),
			SchemaURL: types.StringValue(kd.SchemaURL),
			DependsOn: append([]string{}, kd.DependsOn...),
			Schema:    types.StringValue(string(kd.Schema)),
		})
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &m)...)
}

// dataSourceClient takes the provider's client out of a configure request.
func dataSourceClient(req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) *client.Client {
	if req.ProviderData == nil {
		return nil
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected provider data", "expected *client.Client")
		return nil
	}
	return c
}
