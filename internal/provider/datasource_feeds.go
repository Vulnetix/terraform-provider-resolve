package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/vulnetix/terraform-provider-vulnetix/internal/client"
)

// Both data sources here read the detection rule feed routes
// (vdb-site handler/detection_rule_feeds.go). They are read-only views: feeds
// are not a config kind yet, so they are created and published in the
// console. Feeds need a Pro plan or higher.

type feedModel struct {
	UUID          types.String     `tfsdk:"uuid"`
	Name          types.String     `tfsdk:"name"`
	Description   types.String     `tfsdk:"description"`
	Status        types.String     `tfsdk:"status"`
	FeedType      types.String     `tfsdk:"feed_type"`
	SchemaVersion types.Int64      `tfsdk:"schema_version"`
	RuleCount     types.Int64      `tfsdk:"rule_count"`
	PublicBaseURL types.String     `tfsdk:"public_base_url"`
	ManifestURL   types.String     `tfsdk:"manifest_url"`
	PublishedAt   types.Int64      `tfsdk:"published_at"`
	RevokedAt     types.Int64      `tfsdk:"revoked_at"`
	CreatedAt     types.Int64      `tfsdk:"created_at"`
	UpdatedAt     types.Int64      `tfsdk:"updated_at"`
	Entries       []feedEntryModel `tfsdk:"entries"`
}

type feedEntryModel struct {
	ID         types.Int64  `tfsdk:"id"`
	RuleTable  types.String `tfsdk:"rule_table"`
	RuleType   types.String `tfsdk:"rule_type"`
	RuleUUID   types.String `tfsdk:"rule_uuid"`
	RuleName   types.String `tfsdk:"rule_name"`
	RuleSource types.String `tfsdk:"rule_source"`
	MatchedCve types.String `tfsdk:"matched_cve"`
	Disabled   types.Bool   `tfsdk:"disabled"`
	AddedAt    types.Int64  `tfsdk:"added_at"`
}

func entryAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id":          schema.Int64Attribute{Computed: true, MarkdownDescription: "The entry's id within the feed."},
		"rule_table":  schema.StringAttribute{Computed: true, MarkdownDescription: "Where the rule comes from: `snort_rule`, `yara_rule`, `virtual_patch` or `countermeasure`."},
		"rule_type":   schema.StringAttribute{Computed: true, MarkdownDescription: "The rule's format, such as `SURICATA` or `SIGMA`."},
		"rule_uuid":   schema.StringAttribute{Computed: true, MarkdownDescription: "The rule's id in its table."},
		"rule_name":   schema.StringAttribute{Computed: true, MarkdownDescription: "The rule's name."},
		"rule_source": schema.StringAttribute{Computed: true, MarkdownDescription: "The rule's source, when known."},
		"matched_cve": schema.StringAttribute{Computed: true, MarkdownDescription: "The vulnerability the rule was added for, when known."},
		"disabled":    schema.BoolAttribute{Computed: true, MarkdownDescription: "Whether the entry is left out of the published feed."},
		"added_at":    schema.Int64Attribute{Computed: true, MarkdownDescription: "When the entry was added, in Unix milliseconds."},
	}
}

func feedAttributes(lookup map[string]schema.Attribute) map[string]schema.Attribute {
	attrs := map[string]schema.Attribute{
		"description":     schema.StringAttribute{Computed: true, MarkdownDescription: "The feed's description."},
		"status":          schema.StringAttribute{Computed: true, MarkdownDescription: "`DRAFT`, `PUBLISHED` or `REVOKED`."},
		"feed_type":       schema.StringAttribute{Computed: true, MarkdownDescription: "The output format, such as `SURICATA`, `MODSECURITY`, `AWS_WAF`, `CLOUDFLARE_WAF` or `SIGMA`."},
		"schema_version":  schema.Int64Attribute{Computed: true, MarkdownDescription: "The feed's file format version."},
		"rule_count":      schema.Int64Attribute{Computed: true, MarkdownDescription: "How many entries the feed holds."},
		"public_base_url": schema.StringAttribute{Computed: true, MarkdownDescription: "Where the published files are served from. Null until the feed is published."},
		"manifest_url":    schema.StringAttribute{Computed: true, MarkdownDescription: "The published feed's `manifest.json`. Null until the feed is published."},
		"published_at":    schema.Int64Attribute{Computed: true, MarkdownDescription: "When the feed was last published, in Unix milliseconds."},
		"revoked_at":      schema.Int64Attribute{Computed: true, MarkdownDescription: "When the feed was revoked, in Unix milliseconds."},
		"created_at":      schema.Int64Attribute{Computed: true, MarkdownDescription: "When the feed was created, in Unix milliseconds."},
		"updated_at":      schema.Int64Attribute{Computed: true, MarkdownDescription: "When the feed was last changed, in Unix milliseconds."},
	}
	for k, v := range lookup {
		attrs[k] = v
	}
	return attrs
}

func optInt(p *int64) types.Int64 {
	if p == nil {
		return types.Int64Null()
	}
	return types.Int64Value(*p)
}

func optString(p *string) types.String {
	if p == nil || *p == "" {
		return types.StringNull()
	}
	return types.StringValue(*p)
}

func feedToModel(f *client.Feed, entries []client.FeedEntry, keep func(client.FeedEntry) bool) feedModel {
	m := feedModel{
		UUID:          types.StringValue(f.UUID),
		Name:          types.StringValue(f.Name),
		Description:   types.StringValue(f.Description),
		Status:        types.StringValue(f.Status),
		FeedType:      types.StringValue(f.FeedType),
		SchemaVersion: types.Int64Value(f.SchemaVersion),
		RuleCount:     types.Int64Value(f.RuleCount),
		PublicBaseURL: optString(f.PublicBaseURL),
		ManifestURL:   types.StringNull(),
		PublishedAt:   optInt(f.PublishedAt),
		RevokedAt:     optInt(f.RevokedAt),
		CreatedAt:     types.Int64Value(f.CreatedAt),
		UpdatedAt:     optInt(f.UpdatedAt),
		Entries:       []feedEntryModel{},
	}
	// The publisher writes manifest.json at the feed's base URL
	// (detection_rule_feeds_publish.go).
	if f.PublicBaseURL != nil && *f.PublicBaseURL != "" {
		m.ManifestURL = types.StringValue(*f.PublicBaseURL + "/manifest.json")
	}
	for _, e := range entries {
		if keep != nil && !keep(e) {
			continue
		}
		m.Entries = append(m.Entries, feedEntryModel{
			ID:         types.Int64Value(e.ID),
			RuleTable:  types.StringValue(e.RuleTable),
			RuleType:   types.StringValue(e.RuleType),
			RuleUUID:   types.StringValue(e.RuleUUID),
			RuleName:   types.StringValue(e.RuleName),
			RuleSource: types.StringValue(e.RuleSource),
			MatchedCve: types.StringValue(e.MatchedCve),
			Disabled:   types.BoolValue(e.Disabled),
			AddedAt:    types.Int64Value(e.AddedAt),
		})
	}
	return m
}

// findFeed resolves a feed by uuid, or by name through the list. Names are
// not unique on the server, so a name matching more than one feed is an error.
func findFeed(ctx context.Context, c *client.Client, uuid, name string) (*client.Feed, []client.FeedEntry, error) {
	if uuid == "" {
		feeds, err := c.ListFeeds(ctx)
		if err != nil {
			return nil, nil, err
		}
		var match []client.Feed
		for _, f := range feeds {
			if f.Name == name {
				match = append(match, f)
			}
		}
		switch len(match) {
		case 0:
			return nil, nil, fmt.Errorf("no detection rule feed is named %q", name)
		case 1:
			uuid = match[0].UUID
		default:
			return nil, nil, fmt.Errorf("%d detection rule feeds are named %q; look it up by uuid", len(match), name)
		}
	}
	f, entries, err := c.GetFeed(ctx, uuid)
	if errors.Is(err, client.ErrNotFound) {
		return nil, nil, fmt.Errorf("detection rule feed %s not found", uuid)
	}
	return f, entries, err
}

// vulnetix_detection_rule_feed

type detectionRuleFeedDataSource struct{ client *client.Client }

func newDetectionRuleFeedDataSource() datasource.DataSource { return &detectionRuleFeedDataSource{} }

var _ datasource.DataSourceWithConfigValidators = (*detectionRuleFeedDataSource)(nil)

func (d *detectionRuleFeedDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_detection_rule_feed"
}

func (d *detectionRuleFeedDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "One detection rule feed and its entries, looked up by `uuid` or `name`. " +
			"Point an IDS, IPS, WAF or SIEM at `manifest_url` to consume it. Feeds need a Pro plan or higher.",
		Attributes: feedAttributes(map[string]schema.Attribute{
			"uuid": schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "The feed's id. Set this or `name`."},
			"name": schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "The feed's name. Set this or `uuid`; it must match exactly one feed."},
			"entries": schema.ListNestedAttribute{
				Computed: true, MarkdownDescription: "Every entry in the feed.",
				NestedObject: schema.NestedAttributeObject{Attributes: entryAttributes()},
			},
		}),
	}
}

func (d *detectionRuleFeedDataSource) ConfigValidators(context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.ExactlyOneOf(path.MatchRoot("uuid"), path.MatchRoot("name")),
	}
}

func (d *detectionRuleFeedDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = dataSourceClient(req, resp)
}

func (d *detectionRuleFeedDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError("Provider not configured", "The vulnetix provider was not configured.")
		return
	}
	var cfg feedModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &cfg)...)
	if resp.Diagnostics.HasError() {
		return
	}
	f, entries, err := findFeed(ctx, d.client, cfg.UUID.ValueString(), cfg.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Could not read the detection rule feed", err.Error())
		return
	}
	m := feedToModel(f, entries, nil)
	resp.Diagnostics.Append(resp.State.Set(ctx, &m)...)
}

// vulnetix_virtual_patches

type virtualPatchesDataSource struct{ client *client.Client }

func newVirtualPatchesDataSource() datasource.DataSource { return &virtualPatchesDataSource{} }

type virtualPatchesModel struct {
	FeedUUID      types.String     `tfsdk:"feed_uuid"`
	FeedName      types.String     `tfsdk:"feed_name"`
	FeedType      types.String     `tfsdk:"feed_type"`
	Status        types.String     `tfsdk:"status"`
	PublicBaseURL types.String     `tfsdk:"public_base_url"`
	ManifestURL   types.String     `tfsdk:"manifest_url"`
	IncludeOff    types.Bool       `tfsdk:"include_disabled"`
	Patches       []feedEntryModel `tfsdk:"patches"`
}

var _ datasource.DataSourceWithConfigValidators = (*virtualPatchesDataSource)(nil)

func (d *virtualPatchesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_patches"
}

func (d *virtualPatchesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The Vulnetix generated virtual patches in one detection rule feed: its entries whose " +
			"`rule_table` is `virtual_patch`. Third-party catalogue rules and countermeasures in the same feed are left out. " +
			"Feeds need a Pro plan or higher.",
		Attributes: map[string]schema.Attribute{
			"feed_uuid": schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "The feed's id. Set this or `feed_name`."},
			"feed_name": schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "The feed's name. Set this or `feed_uuid`; it must match exactly one feed."},
			"include_disabled": schema.BoolAttribute{
				Optional: true, MarkdownDescription: "Also return entries that are disabled, and so left out of the published feed. Defaults to false.",
			},
			"feed_type":       schema.StringAttribute{Computed: true, MarkdownDescription: "The feed's output format, such as `MODSECURITY` or `AWS_WAF`."},
			"status":          schema.StringAttribute{Computed: true, MarkdownDescription: "`DRAFT`, `PUBLISHED` or `REVOKED`."},
			"public_base_url": schema.StringAttribute{Computed: true, MarkdownDescription: "Where the published files are served from."},
			"manifest_url":    schema.StringAttribute{Computed: true, MarkdownDescription: "The published feed's `manifest.json`."},
			"patches": schema.ListNestedAttribute{
				Computed: true, MarkdownDescription: "The virtual patch entries.",
				NestedObject: schema.NestedAttributeObject{Attributes: entryAttributes()},
			},
		},
	}
}

func (d *virtualPatchesDataSource) ConfigValidators(context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.ExactlyOneOf(path.MatchRoot("feed_uuid"), path.MatchRoot("feed_name")),
	}
}

func (d *virtualPatchesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = dataSourceClient(req, resp)
}

func (d *virtualPatchesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError("Provider not configured", "The vulnetix provider was not configured.")
		return
	}
	var cfg virtualPatchesModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &cfg)...)
	if resp.Diagnostics.HasError() {
		return
	}
	f, entries, err := findFeed(ctx, d.client, cfg.FeedUUID.ValueString(), cfg.FeedName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Could not read the detection rule feed", err.Error())
		return
	}
	includeOff := cfg.IncludeOff.ValueBool()
	fm := feedToModel(f, entries, func(e client.FeedEntry) bool {
		return e.RuleTable == client.RuleTableVirtualPatch && (includeOff || !e.Disabled)
	})
	m := virtualPatchesModel{
		FeedUUID:      fm.UUID,
		FeedName:      fm.Name,
		FeedType:      fm.FeedType,
		Status:        fm.Status,
		PublicBaseURL: fm.PublicBaseURL,
		ManifestURL:   fm.ManifestURL,
		IncludeOff:    cfg.IncludeOff,
		Patches:       fm.Entries,
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &m)...)
}
