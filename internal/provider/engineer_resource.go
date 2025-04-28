// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"terraform-provider-devops-bootcamp/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &EngineerResource{}
var _ resource.ResourceWithConfigure = &EngineerResource{}
var _ resource.ResourceWithImportState = &EngineerResource{}

func NewEngineerResource() resource.Resource {
	return &EngineerResource{}
}

// EngineerResource defines the resource implementation.
type EngineerResource struct {
	client *client.Client
}

// EngineerResourceModel describes the resource data model.
type EngineerResourceModel struct {
	Id    types.String `tfsdk:"id"`
	Name  types.String `tfsdk:"name"`
	Email types.String `tfsdk:"email"`
}

func (r *EngineerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_engineer_resource"
}

func (r *EngineerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Engineer resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Engineer identifier",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Engineer name",
				Required:            true,
			},
			"email": schema.StringAttribute{
				MarkdownDescription: "Engineer email",
				Required:            true,
			},
		},
	}
}

func (r *EngineerResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *EngineerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data EngineerResourceModel
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check if the client is configured
	if r.client == nil {
		resp.Diagnostics.AddError(
			"Client Not Configured",
			"The provider client has not been configured. Please configure the provider before using this resource.",
		)
		return
	}

	// Create the engineer using the client
	engineer := client.Engineer{
		Name:  data.Name.ValueString(),
		Email: data.Email.ValueString(),
		// Don't include ID in the create request, let the API assign it
	}

	err := r.client.CreateEngineer(engineer)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Engineer",
			err.Error(),
		)
		return
	}

	// Get engineers from the API to retrieve the ID of the newly created engineer
	clientEngineers, err := r.client.GetEngineers()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Bootcamp Engineers",
			err.Error(),
		)
		return
	}

	// Find the newly created engineer by matching name and email
	for _, eng := range clientEngineers {
		if eng.Name == data.Name.ValueString() && eng.Email == data.Email.ValueString() {
			data.Id = types.StringValue(eng.ID)
			break
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *EngineerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state EngineerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check if the client is configured
	if r.client == nil {
		resp.Diagnostics.AddError(
			"Client Not Configured",
			"The provider client has not been configured. Please configure the provider before using this resource.",
		)
		return
	}

	// Get engineers from the API
	clientEngineers, err := r.client.GetEngineers()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Bootcamp Engineers",
			err.Error(),
		)
		return
	}

	// Find the engineer with the matching ID
	found := false
	for _, eng := range clientEngineers {
		if eng.ID == state.Id.ValueString() {
			state.Name = types.StringValue(eng.Name)
			state.Email = types.StringValue(eng.Email)
			found = true
			break
		}
	}

	// If the engineer is not found, remove it from state
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}

	tflog.Trace(ctx, "read engineer resource")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *EngineerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data EngineerResourceModel
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check if the client is configured
	if r.client == nil {
		resp.Diagnostics.AddError(
			"Client Not Configured",
			"The provider client has not been configured. Please configure the provider before using this resource.",
		)
		return
	}

	// Update the engineer using the client
	engineer := client.Engineer{
		ID:    data.Id.ValueString(),
		Name:  data.Name.ValueString(),
		Email: data.Email.ValueString(),
	}

	err := r.client.UpdateEngineer(engineer)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Update Engineer",
			err.Error(),
		)
		return
	}

	clientEngineer, err := r.client.GetEngineerByID(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Bootcamp Engineer",
			err.Error(),
		)
		return
	}

	data.Name = types.StringValue(clientEngineer.Name)
	data.Email = types.StringValue(clientEngineer.Email)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *EngineerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data EngineerResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Check if the client is configured
	if r.client == nil {
		resp.Diagnostics.AddError(
			"Client Not Configured",
			"The provider client has not been configured. Please configure the provider before using this resource.",
		)
		return
	}

	// Delete the engineer using the client
	err := r.client.DeleteEngineer(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Delete Engineer",
			err.Error(),
		)
		return
	}

	tflog.Trace(ctx, "deleted engineer resource")
}

func (r *EngineerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
