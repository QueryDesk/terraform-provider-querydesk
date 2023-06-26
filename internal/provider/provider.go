// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"os"
	"terraform-provider-querydesk/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure QueryDeskProvider satisfies various provider interfaces.
var _ provider.Provider = &QueryDeskProvider{}

// QueryDeskProvider defines the provider implementation.
type QueryDeskProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version    string
	testClient client.GraphQLClient
}

// QueryDeskProviderModel describes the provider data model.
type QueryDeskProviderModel struct {
	Host   types.String `tfsdk:"host"`
	ApiKey types.String `tfsdk:"api_key"`
}

func (p *QueryDeskProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "querydesk"
	resp.Version = p.version
}

func (p *QueryDeskProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Required: true,
			},
			"api_key": schema.StringAttribute{
				Required:  true,
				Sensitive: true,
			},
		},
	}
}

func (p *QueryDeskProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data QueryDeskProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown QueryDesk API Host",
			"The provider cannot create the QueryDesk API client as there is an unknown configuration value for the QueryDesk API host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the QUERYDESK_HOST environment variable.",
		)
	}

	if data.ApiKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Unknown QueryDesk API Key",
			"The provider cannot create the QueryDesk API client as there is an unknown configuration value for the QueryDesk API key. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the QUERYDESK_API_KEY environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	host := os.Getenv("QUERYDESK_HOST")
	api_key := os.Getenv("QUERYDESK_API_KEY")

	if !data.Host.IsNull() {
		host = data.Host.ValueString()
	}

	if !data.ApiKey.IsNull() {
		api_key = data.ApiKey.ValueString()
	}

	if api_key == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Missing QueryDesk API Key",
			"The provider cannot create the QueryDesk API client as there is a missing or empty value for the QueryDesk API key. "+
				"Set the api key value in the configuration or use the QUERYDESK_API_KEY environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	graphqlClient, err := client.NewClient(&host, &api_key)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create QueryDesk API Client",
			"An unexpected error occurred when creating the QueryDesk API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"QueryDesk Client Error: "+err.Error(),
		)
		return
	}

	var myclient client.GraphQLClient

	myclient = client.GraphQLReq{Client: *graphqlClient}

	if p.testClient != nil {
		myclient = p.testClient
	}

	resp.DataSourceData = myclient
	resp.ResourceData = myclient
}

func (p *QueryDeskProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewDatabaseResource,
	}
}

func (p *QueryDeskProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func New(version string, client client.GraphQLClient) func() provider.Provider {
	return func() provider.Provider {
		return &QueryDeskProvider{
			version:    version,
			testClient: client,
		}
	}
}
