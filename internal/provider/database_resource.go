// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"terraform-provider-querydesk/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &DatabaseResource{}
var _ resource.ResourceWithConfigure = &DatabaseResource{}
var _ resource.ResourceWithImportState = &DatabaseResource{}

func NewDatabaseResource() resource.Resource {
	return &DatabaseResource{}
}

// DatabaseResource defines the resource implementation.
type DatabaseResource struct {
	graphqlClient client.GraphQLClient
}

// DatabaseResourceModel describes the resource data model.
type DatabaseResourceModel struct {
	Id             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	Adapter        types.String `tfsdk:"adapter"`
	Hostname       types.String `tfsdk:"hostname"`
	Database       types.String `tfsdk:"database"`
	Ssl            types.Bool   `tfsdk:"ssl"`
	CaCertFile     types.String `tfsdk:"cacertfile"`
	KeyFile        types.String `tfsdk:"keyfile"`
	CertFile       types.String `tfsdk:"certfile"`
	RestrictAccess types.Bool   `tfsdk:"restrict_access"`
}

func (r *DatabaseResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_database"
}

func (r *DatabaseResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Database resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Database id.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name for users to use to identity the database.",
				Required:            true,
			},
			"adapter": schema.StringAttribute{
				MarkdownDescription: "The adapter to use to establish the connection. Currently only `POSTGRES` and `MYSQL` are supported, but  sql server is on the roadmap.",
				Required:            true,
			},
			"database": schema.StringAttribute{
				MarkdownDescription: "The name of the database to connect to.",
				Required:            true,
			},
			"hostname": schema.StringAttribute{
				MarkdownDescription: "The hostname for connecting to the database, either an ip or url.",
				Required:            true,
			},
			"ssl": schema.BoolAttribute{
				MarkdownDescription: "Set to `true` to turn on ssl connections for this database.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"cacertfile": schema.StringAttribute{
				MarkdownDescription: "The server ca cert to use with ssl connections, `ssl` must be set to `true`.",
				Optional:            true,
				Sensitive:           true,
			},
			"keyfile": schema.StringAttribute{
				MarkdownDescription: "The client key to use with ssl connections, `ssl` must be set to `true`.",
				Optional:            true,
				Sensitive:           true,
			},
			"certfile": schema.StringAttribute{
				MarkdownDescription: "The client cert to use with ssl connections, `ssl` must be set to `true`.",
				Optional:            true,
				Sensitive:           true,
			},
			"restrict_access": schema.BoolAttribute{
				MarkdownDescription: "Whether access to this databases should be explicitly granted to users or if any authenticated user can access it.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
		},
	}
}

func (r *DatabaseResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	graphqlClient, ok := req.ProviderData.(client.GraphQLClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.GraphQLReq, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.graphqlClient = graphqlClient
}

func (r *DatabaseResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *DatabaseResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var adapter client.DatabaseAdapter
	switch data.Adapter.ValueString() {
	case "POSTGRES":
		adapter = client.DatabaseAdapterPostgres
	case "MYSQL":
		adapter = client.DatabaseAdapterMysql
	default:
		resp.Diagnostics.AddError("Unexpected Database Adapter", fmt.Sprintf("Expected `POSTGRES` or `MYSQL`, got: %s.", data.Adapter.String()))
		return
	}

	input := client.CreateDatabaseInput{
		Name:           data.Name.ValueString(),
		Adapter:        adapter,
		Hostname:       data.Hostname.ValueString(),
		Database:       data.Database.ValueString(),
		Ssl:            data.Ssl.ValueBool(),
		Cacertfile:     data.CaCertFile.ValueString(),
		Keyfile:        data.KeyFile.ValueString(),
		Certfile:       data.CertFile.ValueString(),
		RestrictAccess: data.RestrictAccess.ValueBool(),
	}

	graphqlResp, err := r.graphqlClient.CreateDatabase(ctx, input)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating database",
			"Could not create database, unexpected error: "+err.Error(),
		)
		return
	}

	if len(graphqlResp.CreateDatabase.Errors) > 0 {
		resp.Diagnostics.AddError(
			"Error creating database",
			"Could not create database: "+graphqlResp.CreateDatabase.Errors[0].Message,
		)
		return
	}

	data.Id = types.StringValue(graphqlResp.CreateDatabase.Result.Id)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DatabaseResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *DatabaseResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	graphqlResp, err := r.graphqlClient.GetDatabase(ctx, data.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Refresh Resource",
			err.Error(),
		)
		return
	}

	// If id is empty, the resource no longer exists
	if graphqlResp.Database.Id == "" {
		resp.State.RemoveResource(ctx)
		return
	}

	data.Name = types.StringValue(graphqlResp.Database.Name)
	data.Adapter = types.StringValue(string(graphqlResp.Database.Adapter))
	data.Hostname = types.StringValue(graphqlResp.Database.Hostname)
	data.Database = types.StringValue(graphqlResp.Database.Database)
	data.Ssl = types.BoolValue(graphqlResp.Database.Ssl)
	data.RestrictAccess = types.BoolValue(graphqlResp.Database.RestrictAccess)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DatabaseResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *DatabaseResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var adapter client.DatabaseAdapter
	switch data.Adapter.ValueString() {
	case "POSTGRES":
		adapter = client.DatabaseAdapterPostgres
	case "MYSQL":
		adapter = client.DatabaseAdapterMysql
	default:
		resp.Diagnostics.AddError("Unexpected Database Adapter", fmt.Sprintf("Expected `POSTGRES` or `MYSQL`, got: %s.", data.Adapter.String()))
		return
	}

	input := client.UpdateDatabaseInput{
		Name:           data.Name.ValueString(),
		Adapter:        adapter,
		Hostname:       data.Hostname.ValueString(),
		Database:       data.Database.ValueString(),
		Ssl:            data.Ssl.ValueBool(),
		NewCacertfile:  data.CaCertFile.ValueString(),
		NewKeyfile:     data.KeyFile.ValueString(),
		NewCertfile:    data.CertFile.ValueString(),
		RestrictAccess: data.RestrictAccess.ValueBool(),
	}

	graphqlResp, err := r.graphqlClient.UpdateDatabase(ctx, data.Id.ValueString(), input)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating database",
			"Could not update database, unexpected error: "+err.Error(),
		)

		return
	}

	if len(graphqlResp.UpdateDatabase.Errors) > 0 {
		resp.Diagnostics.AddError(
			"Error updating database",
			"Could not update database, unexpected error: "+graphqlResp.UpdateDatabase.Errors[0].Message,
		)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DatabaseResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *DatabaseResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	graphqlResp, err := r.graphqlClient.DeleteDatabase(ctx, data.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to delete database, got error: %s", err),
		)

		return
	}

	if len(graphqlResp.DeleteDatabase.Errors) > 0 {
		resp.Diagnostics.AddError(
			"Error deleting database",
			"Could not delete database, unexpected error: "+graphqlResp.DeleteDatabase.Errors[0].Message,
		)
		return
	}
}

func (r *DatabaseResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
