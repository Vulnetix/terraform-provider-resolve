package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/vulnetix/terraform-provider-resolve/internal/client"
)

// vulnetix_github_installation_repos reads the repositories the tenant's
// GitHub App installations reach (GET /integrations/github/installations/repos),
// from Vulnetix's cache of them, so a plan never calls GitHub.

type githubReposDataSource struct{ client *client.Client }

func newGitHubInstallationReposDataSource() datasource.DataSource { return &githubReposDataSource{} }

type githubReposModel struct {
	InstallationID types.Int64       `tfsdk:"installation_id"`
	Query          types.String      `tfsdk:"query"`
	IncludeArchive types.Bool        `tfsdk:"include_archived"`
	Repos          []githubRepoModel `tfsdk:"repos"`
	FullNames      []string          `tfsdk:"full_names"`
}

type githubRepoModel struct {
	InstallationID types.Int64  `tfsdk:"installation_id"`
	AccountLogin   types.String `tfsdk:"account_login"`
	RepoID         types.Int64  `tfsdk:"repo_id"`
	OwnerLogin     types.String `tfsdk:"owner_login"`
	Name           types.String `tfsdk:"name"`
	FullName       types.String `tfsdk:"full_name"`
	Private        types.Bool   `tfsdk:"private"`
	Archived       types.Bool   `tfsdk:"archived"`
	Host           types.String `tfsdk:"host"`
	UpdatedAt      types.Int64  `tfsdk:"updated_at"`
}

func (d *githubReposDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_github_installation_repos"
}

func (d *githubReposDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The repositories the tenant's active GitHub App installations reach " +
			"(`GET /integrations/github/installations/repos`), ordered by full name. It reads Vulnetix's cache, which " +
			"GitHub webhooks keep current, so a plan never calls GitHub. Use `full_names` to scope a " +
			"`vulnetix_vex_source` or a `vulnetix_jail_policy` to every repository the App can see.",
		Attributes: map[string]schema.Attribute{
			"installation_id": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "Only this installation's repositories.",
				Validators:          []validator.Int64{int64validator.AtLeast(1)},
			},
			"query":            schema.StringAttribute{Optional: true, MarkdownDescription: "Only repositories whose `owner/name` contains this, ignoring case."},
			"include_archived": schema.BoolAttribute{Optional: true, MarkdownDescription: "Keep archived repositories. Defaults to `false`."},
			"full_names": schema.ListAttribute{
				Computed: true, ElementType: types.StringType,
				MarkdownDescription: "The repositories' `owner/name`, in order.",
			},
			"repos": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The matching repositories.",
				NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
					"installation_id": schema.Int64Attribute{Computed: true, MarkdownDescription: "The GitHub App installation that reaches it."},
					"account_login":   schema.StringAttribute{Computed: true, MarkdownDescription: "The account the App is installed on."},
					"repo_id":         schema.Int64Attribute{Computed: true, MarkdownDescription: "GitHub's repository id."},
					"owner_login":     schema.StringAttribute{Computed: true, MarkdownDescription: "The repository's owner."},
					"name":            schema.StringAttribute{Computed: true, MarkdownDescription: "The repository's name."},
					"full_name":       schema.StringAttribute{Computed: true, MarkdownDescription: "`owner/name`."},
					"private":         schema.BoolAttribute{Computed: true, MarkdownDescription: "Whether the repository is private."},
					"archived":        schema.BoolAttribute{Computed: true, MarkdownDescription: "Whether the repository is archived."},
					"host":            schema.StringAttribute{Computed: true, MarkdownDescription: "The GitHub host, such as `github.com`."},
					"updated_at":      schema.Int64Attribute{Computed: true, MarkdownDescription: "When Vulnetix last refreshed the row, in Unix milliseconds."},
				}},
			},
		},
	}
}

func (d *githubReposDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = dataSourceClient(req, resp)
}

func (d *githubReposDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError("Provider not configured", "The vulnetix provider was not configured.")
		return
	}
	var m githubReposModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &m)...)
	if resp.Diagnostics.HasError() {
		return
	}
	repos, err := d.client.ListGitHubInstallationRepos(ctx, client.GitHubRepoFilter{
		InstallationID: m.InstallationID.ValueInt64(),
		Query:          m.Query.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Could not list GitHub installation repositories", err.Error())
		return
	}
	keepArchived := m.IncludeArchive.ValueBool()
	m.Repos, m.FullNames = []githubRepoModel{}, []string{}
	for _, r := range repos {
		if r.Archived && !keepArchived {
			continue
		}
		m.FullNames = append(m.FullNames, r.FullName)
		m.Repos = append(m.Repos, githubRepoModel{
			InstallationID: types.Int64Value(r.InstallationID),
			AccountLogin:   types.StringValue(r.AccountLogin),
			RepoID:         types.Int64Value(r.RepoID),
			OwnerLogin:     types.StringValue(r.OwnerLogin),
			Name:           types.StringValue(r.Name),
			FullName:       types.StringValue(r.FullName),
			Private:        types.BoolValue(r.Private),
			Archived:       types.BoolValue(r.Archived),
			Host:           types.StringValue(r.Host),
			UpdatedAt:      types.Int64Value(r.UpdatedAt),
		})
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &m)...)
}
