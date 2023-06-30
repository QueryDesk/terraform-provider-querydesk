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
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &DatabaseUserResource{}
var _ resource.ResourceWithConfigure = &DatabaseUserResource{}
var _ resource.ResourceWithImportState = &DatabaseUserResource{}

func NewDatabaseUserResource() resource.Resource {
	return &DatabaseUserResource{}
}

// DatabaseUserResource defines the resource implementation.
type DatabaseUserResource struct {
	graphqlClient client.GraphQLClient
}

// DatabaseUserResourceModel describes the resource data model.
type DatabaseUserResourceModel struct {
	Id              types.String `tfsdk:"id"`
	DatabaseId      types.String `tfsdk:"database_id"`
	Description     types.String `tfsdk:"description"`
	Username        types.String `tfsdk:"username"`
	Password        types.String `tfsdk:"password"`
	ReviewsRequired types.Int64  `tfsdk:"reviews_required"`
}

func (r *DatabaseUserResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_database_user"
}

func (r *DatabaseUserResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "DatabaseUser resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Info shown in the UI to help identity available users.",
				Optional:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "The user to authenticate with.",
				Required:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "The password to authenticate the user with.",
				Required:            true,
				Sensitive:           true,
			},
			"reviews_required": schema.Int64Attribute{
				MarkdownDescription: "How many reviews are required to use this user. Can be set to 0 to not require reviews.",
				Required:            true,
			},
			"database_id": schema.StringAttribute{
				MarkdownDescription: "Identifier of the related database.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *DatabaseUserResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *DatabaseUserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *DatabaseUserResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	input := client.CreateCredentialInput{
		DatabaseId:      data.DatabaseId.ValueString(),
		Description:     data.Description.ValueString(),
		Password:        data.Password.ValueString(),
		ReviewsRequired: int(data.ReviewsRequired.ValueInt64()),
		Username:        data.Username.ValueString(),
	}

	graphqlResp, err := r.graphqlClient.CreateCredential(ctx, input)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating database user",
			"Could not create database user, unexpected error: "+err.Error(),
		)
		return
	}

	if len(graphqlResp.CreateCredential.Errors) > 0 {
		resp.Diagnostics.AddError(
			"Error creating database user",
			"Could not create database user: "+graphqlResp.CreateCredential.Errors[0].Message,
		)
		return
	}

	data.Id = types.StringValue(graphqlResp.CreateCredential.Result.Id)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DatabaseUserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *DatabaseUserResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	graphqlResp, err := r.graphqlClient.GetCredential(ctx, data.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Refresh Resource",
			err.Error(),
		)
		return
	}

	// If id is empty, the resource no longer exists
	if graphqlResp.Credential.Id == "" {
		resp.State.RemoveResource(ctx)
		return
	}

	if graphqlResp.Credential.Description != "" {
		data.Description = types.StringValue(graphqlResp.Credential.Description)
	}

	data.Username = types.StringValue(graphqlResp.Credential.Username)
	data.ReviewsRequired = types.Int64Value(int64(graphqlResp.Credential.ReviewsRequired))
	data.DatabaseId = types.StringValue(graphqlResp.Credential.Database.Id)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DatabaseUserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *DatabaseUserResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	input := client.UpdateCredentialInput{
		Description:     data.Description.ValueString(),
		NewPassword:     data.Password.ValueString(),
		ReviewsRequired: int(data.ReviewsRequired.ValueInt64()),
		Username:        data.Username.ValueString(),
	}

	graphqlResp, err := r.graphqlClient.UpdateCredential(ctx, data.Id.ValueString(), input)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating database user",
			"Could not update database user, unexpected error: "+err.Error(),
		)

		return
	}

	if len(graphqlResp.UpdateCredential.Errors) > 0 {
		resp.Diagnostics.AddError(
			"Error updating database user",
			"Could not update database user, unexpected error: "+graphqlResp.UpdateCredential.Errors[0].Message,
		)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DatabaseUserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *DatabaseUserResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	graphqlResp, err := r.graphqlClient.DeleteCredential(ctx, data.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to delete database user, got error: %s", err),
		)

		return
	}

	if len(graphqlResp.DeleteCredential.Errors) > 0 {
		resp.Diagnostics.AddError(
			"Error deleting database user",
			"Could not delete database user, unexpected error: "+graphqlResp.DeleteCredential.Errors[0].Message,
		)
		return
	}
}

func (r *DatabaseUserResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
