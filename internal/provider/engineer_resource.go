// // Copyright (c) HashiCorp, Inc.
// // SPDX-License-Identifier: MPL-2.0

package provider

// import (
// 	"context"
// 	"fmt"

// 	"terraform-provider-devops-bootcamp/internal/client"

// 	"github.com/hashicorp/terraform-plugin-framework/path"
// 	"github.com/hashicorp/terraform-plugin-framework/resource"
// 	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
// 	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
// 	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
// 	"github.com/hashicorp/terraform-plugin-log/tflog"
// )

// // Ensure provider defined types fully satisfy framework interfaces.
// var _ resource.Resource = &EngineerResource{}
// var _ resource.ResourceWithConfigure = &EngineerResource{}
// var _ resource.ResourceWithImportState = &EngineerResource{}

// func NewEngineerResource() resource.Resource {
// 	return &EngineerResource{}
// }

// // EngineerResource defines the resource implementation.
// type EngineerResource struct {
// 	client *client.Client
// }

// // EngineerResourceModel describes the resource data model.
// // type EngineerResourceModel struct {
// // 	Engineers []EngineerDataSourceModel `tfsdk:"engineers"`
// // }

// type EngineerResourceModel struct {
// 	Id    string `tfsdk:"id"`
// 	Name  string `tfsdk:"name"`
// 	Email string `tfsdk:"email"`
// }

// func (r *EngineerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
// 	resp.TypeName = req.ProviderTypeName + "_engineer"
// }

// func (r *EngineerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
// 	resp.Schema = schema.Schema{
// 		Attributes: map[string]schema.Attribute{
// 			"id": schema.StringAttribute{
// 				MarkdownDescription: "Engineer identifier",
// 				Computed:            true,
// 				PlanModifiers: []planmodifier.String{
// 					stringplanmodifier.UseStateForUnknown(),
// 				},
// 			},
// 			"name": schema.StringAttribute{
// 				MarkdownDescription: "Engineer name",
// 				Required:            true,
// 			},
// 			"email": schema.StringAttribute{
// 				MarkdownDescription: "Engineer email",
// 				Required:            true,
// 			},
// 		},
// 	}
// }

// func (r *EngineerResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
// 	if req.ProviderData == nil {
// 		return
// 	}

// 	client, ok := req.ProviderData.(*client.Client)

// 	if !ok {
// 		resp.Diagnostics.AddError(
// 			"Unexpected Resource Configure Type",
// 			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
// 		)

// 		return
// 	}

// 	r.client = client
// }

// func (r *EngineerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
// 	var data EngineerModel
// 	diags := req.Plan.Get(ctx, &data)
// 	resp.Diagnostics.Append(diags...)
// 	if resp.Diagnostics.HasError() {
// 		return
// 	}

// 	// Check if the client is configured
// 	if r.client == nil {
// 		resp.Diagnostics.AddError(
// 			"Client Not Configured",
// 			"The provider client has not been configured. Please configure the provider before using this resource.",
// 		)
// 		return
// 	}

// 	// Create the engineer using the client
// 	engineer := client.Engineer{
// 		Name:  data.Name,
// 		Email: data.Email,
// 	}

// 	err := r.client.CreateEngineer(engineer)
// 	if err != nil {
// 		resp.Diagnostics.AddError(
// 			"Unable to Create Engineer",
// 			err.Error(),
// 		)
// 		return
// 	}

// 	// Get engineers from the API to retrieve the ID of the newly created engineer
// 	clientEngineers, err := r.client.GetEngineers()
// 	if err != nil {
// 		resp.Diagnostics.AddError(
// 			"Unable to Read Bootcamp Engineers",
// 			err.Error(),
// 		)
// 		return
// 	}

// 	// Find the newly created engineer by matching name and email
// 	for _, eng := range clientEngineers {
// 		if eng.Name == data.Name && eng.Email == data.Email {
// 			data.Id = eng.ID
// 			break
// 		}
// 	}

// 	tflog.Trace(ctx, "created a resource")

// 	// Save data into Terraform state
// 	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
// }

// func (r *EngineerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
// 	var data EngineerModel

// 	// Read Terraform prior state data into the model
// 	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

// 	if resp.Diagnostics.HasError() {
// 		return
// 	}

// 	// Check if the client is configured
// 	if r.client == nil {
// 		resp.Diagnostics.AddError(
// 			"Client Not Configured",
// 			"The provider client has not been configured. Please configure the provider before using this resource.",
// 		)
// 		return
// 	}

// 	// Get engineers from the API
// 	clientEngineers, err := r.client.GetEngineers()
// 	if err != nil {
// 		resp.Diagnostics.AddError(
// 			"Unable to Read Bootcamp Engineers",
// 			err.Error(),
// 		)
// 		return
// 	}

// 	// Find the engineer with the matching ID
// 	found := false
// 	for _, eng := range clientEngineers {
// 		if eng.ID == data.Id {
// 			data.Name = eng.Name
// 			data.Email = eng.Email
// 			found = true
// 			break
// 		}
// 	}

// 	// If the engineer is not found, remove it from state
// 	if !found {
// 		resp.State.RemoveResource(ctx)
// 		return
// 	}

// 	tflog.Trace(ctx, "read engineer resource")

// 	// Save updated data into Terraform state
// 	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
// }

// func (r *EngineerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
// 	var data EngineerModel

// 	// Read Terraform plan data into the model
// 	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

// 	if resp.Diagnostics.HasError() {
// 		return
// 	}

// 	// Check if the client is configured
// 	if r.client == nil {
// 		resp.Diagnostics.AddError(
// 			"Client Not Configured",
// 			"The provider client has not been configured. Please configure the provider before using this resource.",
// 		)
// 		return
// 	}

// 	// Update the engineer using the client
// 	engineer := client.Engineer{
// 		ID:    data.Id,
// 		Name:  data.Name,
// 		Email: data.Email,
// 	}

// 	err := r.client.UpdateEngineer(engineer)
// 	if err != nil {
// 		resp.Diagnostics.AddError(
// 			"Unable to Update Engineer",
// 			err.Error(),
// 		)
// 		return
// 	}

// 	tflog.Trace(ctx, "updated engineer resource")

// 	// Save updated data into Terraform state
// 	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
// }

// func (r *EngineerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
// 	var data EngineerModel

// 	// Read Terraform prior state data into the model
// 	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

// 	if resp.Diagnostics.HasError() {
// 		return
// 	}

// 	// Check if the client is configured
// 	if r.client == nil {
// 		resp.Diagnostics.AddError(
// 			"Client Not Configured",
// 			"The provider client has not been configured. Please configure the provider before using this resource.",
// 		)
// 		return
// 	}

// 	// Delete the engineer using the client
// 	err := r.client.DeleteEngineer(data.Id)
// 	if err != nil {
// 		resp.Diagnostics.AddError(
// 			"Unable to Delete Engineer",
// 			err.Error(),
// 		)
// 		return
// 	}

// 	tflog.Trace(ctx, "deleted engineer resource")
// }

// func (r *EngineerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
// 	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
// }
