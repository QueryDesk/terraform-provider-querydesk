package provider

import (
	"context"
	"fmt"
	"terraform-provider-querydesk/internal/client"

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
	client *client.Client
}

type DatabaseDataSourceModel struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
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
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *DatabaseDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *DatabaseDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data DatabaseDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	database, err := d.client.GetDatabase(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to get database",
			err.Error(),
		)
		return
	}

	data = DatabaseDataSourceModel{
		Id:   types.StringValue(database.Id),
		Name: types.StringValue(database.Name),
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
