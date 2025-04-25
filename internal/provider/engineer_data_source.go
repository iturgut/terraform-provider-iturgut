package provider

import (
	"context"
	"fmt"

	"terraform-provider-devops-bootcamp/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ datasource.DataSource = &EngineerDataSource{}

func NewEngineerDataSource() datasource.DataSource {
	return &EngineerDataSource{}
}

type EngineerDataSource struct {
	client *client.Client
}

type EngineerDataSourceModel struct {
	Name  types.String `tfsdk:"name"`
	Id    types.String `tfsdk:"id"`
	Email types.String `tfsdk:"email"`
}

func (d *EngineerDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_engineer_data_source"
}

func (d *EngineerDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Engineer data source",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Engineer identifier",
				Optional:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Engineer name",
				Computed:            true,
			},
			"email": schema.StringAttribute{
				MarkdownDescription: "Engineer email",
				Computed:            true,
			},
		},
	}
}

func (d *EngineerDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *EngineerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data EngineerDataSourceModel

	// Check if the client is configured
	if d.client == nil {
		resp.Diagnostics.AddError(
			"Client Not Configured",
			"The provider client has not been configured. Please configure the provider before using this data source.",
		)
		return
	}

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	engineer, err := d.client.GetEngineerByID(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Bootcamp Engineer",
			err.Error(),
		)
		return
	}

	// Map client.Engineer to EngineerDataSourceModel
	engineerState := EngineerDataSourceModel{
		Id:    types.StringValue(engineer.ID),
		Name:  types.StringValue(engineer.Name),
		Email: types.StringValue(engineer.Email),
	}
	data = engineerState

	// Write logs using the tflog package
	tflog.Trace(ctx, "read engineers data source")

	diags := resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Write logs using the tflog package
	tflog.Trace(ctx, "read engineers data source")

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
