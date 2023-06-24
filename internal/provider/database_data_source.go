package provider

import (
	"context"
	"fmt"
	"terraform-provider-querydesk/internal/client"

	"github.com/Khan/genqlient/graphql"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &DatabaseDataSource{}
var _ datasource.DataSourceWithConfigure = &DatabaseDataSource{}

func NewDatabaseDataSource() datasource.DataSource {
	return &DatabaseDataSource{}
}

type DatabaseDataSource struct {
	graphqlClient *graphql.Client
}

type DatabaseDataSourceModel struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	// TODO: make an enum
	Adapter        types.String `tfsdk:"adapter"`
	Hostname       types.String `tfsdk:"hostname"`
	Database       types.String `tfsdk:"database"`
	Ssl            types.Bool   `tfsdk:"ssl"`
	RestrictAccess types.Bool   `tfsdk:"restrict_access"`
}

func (d *DatabaseDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_database"
}

func (d *DatabaseDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Database data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Database id",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the database connection",
				Computed:            true,
			},
			"adapter": schema.StringAttribute{
				MarkdownDescription: "Database name",
				Computed:            true,
			},
			"database": schema.StringAttribute{
				MarkdownDescription: "Database name",
				Computed:            true,
			},
			"hostname": schema.StringAttribute{
				MarkdownDescription: "Database name",
				Computed:            true,
			},
			"ssl": schema.BoolAttribute{
				MarkdownDescription: "Database name",
				Computed:            true,
			},
			"restrict_access": schema.BoolAttribute{
				MarkdownDescription: "Database name",
				Computed:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *DatabaseDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*graphql.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.graphqlClient = client
}

func (d *DatabaseDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data DatabaseDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	remoteData, err := client.GetDatabase(ctx, *d.graphqlClient, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to get database",
			err.Error(),
		)
		return
	}

	data.Name = types.StringValue(remoteData.Database.Name)
	data.Adapter = types.StringValue(remoteData.Database.Adapter)
	data.Hostname = types.StringValue(remoteData.Database.Hostname)
	data.Database = types.StringValue(remoteData.Database.Database)
	data.Ssl = types.BoolValue(remoteData.Database.Ssl)
	data.RestrictAccess = types.BoolValue(remoteData.Database.RestrictAccess)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
