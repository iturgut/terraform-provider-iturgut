package provider

import (
	"context"
	"os"
	"terraform-provider-devops-bootcamp/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure bootcampProvider satisfies various provider interfaces.
var _ provider.Provider = &bootcampProvider{}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &bootcampProvider{
			version: version,
		}
	}
}

// bootcampProviderModel describes the provider data model.
type bootcampProviderModel struct {
	Host types.String `tfsdk:"host"`
}

// bootcampProvider defines the provider implementation.
type bootcampProvider struct {
	version string
}

func (p *bootcampProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "devops-bootcamp"
	resp.Version = p.version
}

func (p *bootcampProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				MarkdownDescription: "The host of the bootcamp API",
				Optional:            true,
			},
		},
	}
}

// retrieves values fron configuration
// checks for unknown config values
// retrieves values from env variables
// creates API client
// stores configured client for data source and resource usage
func (p *bootcampProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Bootcamp client")
	// retrieve config values
	var config bootcampProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown Bootcamp API Host",
			"The provider cannot create the Bootcamp API client as there is an unknown configuration value for the Bootcamp API host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the BOOTCAMP_HOST environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	host := os.Getenv("BOOTCAMP_HOST")

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing Bootcamp API Host",
			"The provider cannot create the Bootcamp API client as there is a missing or empty value for the Bootcamp API host. "+
				"Set the host value in the configuration or use the BOOTCAMP_HOST environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "bootcamp_host", host)
	tflog.Debug(ctx, "Creating Bootcamp client")

	client, err := client.NewClient(&host)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Bootcamp API Client",
			"An unexpected error occurred when creating the bootcamp API client. "+
				"Bootcamp Client Error: "+err.Error(),
		)
		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured Bootcamp client", map[string]any{"success": true})
}

func (p *bootcampProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		//NewEngineerResource,
		//NewDevResource,
	}
}

func (p *bootcampProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewEngineerDataSource,
		//NewDevDataSource,
	}
}
