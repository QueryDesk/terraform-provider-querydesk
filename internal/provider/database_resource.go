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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
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
	client *client.Client
}

// DatabaseResourceModel describes the resource data model.
type DatabaseResourceModel struct {
	Id       types.String `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
	// TODO: make an enum
	Adapter         types.String `tfsdk:"adapter"`
	Hostname        types.String `tfsdk:"hostname"`
	Database        types.String `tfsdk:"database"`
	ReviewsRequired types.Int64  `tfsdk:"reviews_required"`
	Ssl             types.Bool   `tfsdk:"ssl"`
	CaCertFile      types.String `tfsdk:"cacertfile"`
	KeyFile         types.String `tfsdk:"keyfile"`
	CertFile        types.String `tfsdk:"certfile"`
	RestrictAccess  types.Bool   `tfsdk:"restrict_access"`
	AgentId         types.String `tfsdk:"agent_id"`
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
			"username": schema.StringAttribute{
				MarkdownDescription: "Database username",
				Required:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "Database name",
				Required:            true,
				Sensitive:           true,
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
			"reviews_required": schema.Int64Attribute{
				MarkdownDescription: "Database name",
				Optional:            true,
				Computed:            true,
			},
			"ssl": schema.BoolAttribute{
				MarkdownDescription: "Database name",
				Optional:            true,
				Computed:            true,
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
			},
			"agent_id": schema.StringAttribute{
				MarkdownDescription: "Database name",
				Optional:            true,
			},
		},
	}
}

func (r *DatabaseResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *DatabaseResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *DatabaseResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	rb := client.Database{
		Name:            data.Name.ValueString(),
		Username:        data.Username.ValueString(),
		Password:        data.Password.ValueString(),
		Adapter:         data.Adapter.ValueString(),
		Hostname:        data.Hostname.ValueString(),
		Database:        data.Database.ValueString(),
		ReviewsRequired: data.ReviewsRequired.ValueInt64(),
		Ssl:             data.Ssl.ValueBool(),
		CaCertFile:      data.CaCertFile.ValueString(),
		KeyFile:         data.KeyFile.ValueString(),
		CertFile:        data.CertFile.ValueString(),
		RestrictAccess:  data.RestrictAccess.ValueBool(),
		AgentId:         data.AgentId.ValueString(),
	}

	database, err := r.client.CreateDatabase(rb)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating database",
			"Could not create database, unexpected error: "+err.Error(),
		)
		return
	}

	// For the purposes of this example code, hardcoding a response value to
	// save into the Terraform state.
	data.Id = types.StringValue(database.Id)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a resource")

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

	database, err := r.client.GetDatabase(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to get database",
			err.Error(),
		)
		return
	}

	data.Name = types.StringValue(database.Name)
	data.Username = types.StringValue(database.Username)
	data.Password = types.StringValue(database.Password)
	data.Adapter = types.StringValue(database.Adapter)
	data.Hostname = types.StringValue(database.Hostname)
	data.Database = types.StringValue(database.Database)
	data.ReviewsRequired = types.Int64Value(database.ReviewsRequired)
	data.Ssl = types.BoolValue(database.Ssl)
	data.CaCertFile = types.StringValue(database.CaCertFile)
	data.KeyFile = types.StringValue(database.KeyFile)
	data.CertFile = types.StringValue(database.CertFile)
	data.RestrictAccess = types.BoolValue(database.RestrictAccess)
	data.AgentId = types.StringValue(database.AgentId)

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

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update example, got error: %s", err))
	//     return
	// }

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

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete example, got error: %s", err))
	//     return
	// }
}

func (r *DatabaseResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
