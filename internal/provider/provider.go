package provider

import (
	"context"
	"os"

	"github.com/algolia/algoliasearch-client-go/v4/algolia/search"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/ECCOShoes/terraform-provider-algolia/internal/client"
)

// Ensure algoliaProvider satisfies the provider.Provider interface.
var _ provider.Provider = (*algoliaProvider)(nil)

type algoliaProvider struct {
	version string
}

type providerModel struct {
	AppID  types.String `tfsdk:"app_id"`
	APIKey types.String `tfsdk:"api_key"`
}

// New returns a function that constructs the provider, as required by
// providerserver.Serve.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &algoliaProvider{version: version}
	}
}

func (p *algoliaProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "algolia"
	resp.Version = p.version
}

func (p *algoliaProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Terraform provider for Algolia.",
		Attributes: map[string]schema.Attribute{
			"app_id": schema.StringAttribute{
				Optional:    true,
				Description: "Algolia application ID. May also be set via the ALGOLIA_APP_ID environment variable.",
			},
			"api_key": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "Algolia Admin API key with permission to manage index settings and API keys. May also be set via the ALGOLIA_API_KEY environment variable.",
			},
		},
	}
}

func (p *algoliaProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var cfg providerModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &cfg)...)
	if resp.Diagnostics.HasError() {
		return
	}

	appID := valueOrEnv(cfg.AppID, "ALGOLIA_APP_ID")
	apiKey := valueOrEnv(cfg.APIKey, "ALGOLIA_API_KEY")

	requireAttr(resp, appID, "app_id", "ALGOLIA_APP_ID")
	requireAttr(resp, apiKey, "api_key", "ALGOLIA_API_KEY")
	if resp.Diagnostics.HasError() {
		return
	}

	c, err := client.New(client.Config{AppID: appID, APIKey: apiKey})
	if err != nil {
		resp.Diagnostics.AddError("Unable to create Algolia client", err.Error())
		return
	}

	resp.ResourceData = c
}

func (p *algoliaProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewAPIKeyResource,
		NewIndexConfigResource,
	}
}

func (p *algoliaProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return nil
}

func valueOrEnv(v types.String, env string) string {
	if !v.IsNull() && !v.IsUnknown() && v.ValueString() != "" {
		return v.ValueString()
	}
	return os.Getenv(env)
}

func requireAttr(resp *provider.ConfigureResponse, value, attr, env string) {
	if value == "" {
		resp.Diagnostics.AddError(
			"Missing Algolia configuration",
			"The provider requires "+attr+" to be set, either in the provider block or via the "+env+" environment variable.",
		)
	}
}

// clientFromProviderData extracts the Algolia client from resource configure
// data, adding a diagnostic if the type is unexpected.
func clientFromProviderData(providerData any, diags interface {
	AddError(summary, detail string)
}) *search.APIClient {
	if providerData == nil {
		return nil
	}
	c, ok := providerData.(*search.APIClient)
	if !ok {
		diags.AddError(
			"Unexpected provider data type",
			"Expected *search.APIClient from the provider configuration.",
		)
		return nil
	}
	return c
}
