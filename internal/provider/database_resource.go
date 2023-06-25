// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"terraform-provider-querydesk/internal/client"

	"github.com/Khan/genqlient/graphql"
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
	graphqlClient *graphql.Client
}

// DatabaseResourceModel describes the resource data model.
type DatabaseResourceModel struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	// TODO: make an enum
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
				MarkdownDescription: "Database identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Database name",
				Required:            true,
			},
			"adapter": schema.StringAttribute{
				MarkdownDescription: "Database name",
				Required:            true,
			},
			"database": schema.StringAttribute{
				MarkdownDescription: "Database name",
				Required:            true,
			},
			"hostname": schema.StringAttribute{
				MarkdownDescription: "Database name",
				Required:            true,
			},
			"ssl": schema.BoolAttribute{
				MarkdownDescription: "Database name",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"cacertfile": schema.StringAttribute{
				MarkdownDescription: "Database name",
				Optional:            true,
				Sensitive:           true,
			},
			"keyfile": schema.StringAttribute{
				MarkdownDescription: "Database name",
				Optional:            true,
				Sensitive:           true,
			},
			"certfile": schema.StringAttribute{
				MarkdownDescription: "Database name",
				Optional:            true,
				Sensitive:           true,
			},
			"restrict_access": schema.BoolAttribute{
				MarkdownDescription: "Database name",
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

	client, ok := req.ProviderData.(*graphql.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.graphqlClient = client
}

func (r *DatabaseResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *DatabaseResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	input := client.CreateDatabaseInput{
		Name:           data.Name.ValueString(),
		Adapter:        data.Adapter.ValueString(),
		Hostname:       data.Hostname.ValueString(),
		Database:       data.Database.ValueString(),
		Ssl:            data.Ssl.ValueBool(),
		Cacertfile:     data.CaCertFile.ValueString(),
		Keyfile:        data.KeyFile.ValueString(),
		Certfile:       data.CertFile.ValueString(),
		RestrictAccess: data.RestrictAccess.ValueBool(),
	}

	remoteData, err := client.CreateDatabase(ctx, *r.graphqlClient, input)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating database",
			"Could not create database, unexpected error: "+err.Error(),
		)
		return
	}

	// TODO: handle graphql errors

	data.Id = types.StringValue(remoteData.CreateDatabase.Result.Id)
	data.Ssl = types.BoolValue(remoteData.CreateDatabase.Result.Ssl)
	data.RestrictAccess = types.BoolValue(remoteData.CreateDatabase.Result.RestrictAccess)

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

	remoteData, err := client.GetDatabase(ctx, *r.graphqlClient, data.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to get database",
			err.Error(),
		)
		return
	}

	if remoteData.Database.Id == "" {
		resp.Diagnostics.AddError(
			"Unable to get database",
			fmt.Sprintf("Database with id %s not found", data.Id.ValueString()),
		)
		return
	}

	data.Name = types.StringValue(remoteData.Database.Name)
	data.Adapter = types.StringValue(remoteData.Database.Adapter)
	data.Hostname = types.StringValue(remoteData.Database.Hostname)
	data.Database = types.StringValue(remoteData.Database.Database)
	data.Ssl = types.BoolValue(remoteData.Database.Ssl)
	data.RestrictAccess = types.BoolValue(remoteData.Database.RestrictAccess)

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

	input := client.UpdateDatabaseInput{
		Name:           data.Name.ValueString(),
		Adapter:        data.Adapter.ValueString(),
		Hostname:       data.Hostname.ValueString(),
		Database:       data.Database.ValueString(),
		Ssl:            data.Ssl.ValueBool(),
		NewCacertfile:  data.CaCertFile.ValueString(),
		NewKeyfile:     data.KeyFile.ValueString(),
		NewCertfile:    data.CertFile.ValueString(),
		RestrictAccess: data.RestrictAccess.ValueBool(),
	}

	_, err := client.UpdateDatabase(ctx, *r.graphqlClient, data.Id.ValueString(), input)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating database",
			"Could not create database, unexpected error: "+err.Error(),
		)
		return
	}

	// TODO: handle graphql errors

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

	_, err := client.DeleteDatabase(ctx, *r.graphqlClient, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete database, got error: %s", err))
		return
	}
	// TODO: handle graphql errors
}

func (r *DatabaseResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
