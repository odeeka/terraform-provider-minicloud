package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/odeeka/minicloud-client-go/minicloud"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &minicloudProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &minicloudProvider{
			version: version,
		}
	}
}

// minicloudProvider is the provider implementation.
type minicloudProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// minicloudProviderModel maps provider schema data to a Go type.
type minicloudProviderModel struct {
	Host     types.String `tfsdk:"host"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

// Metadata returns the provider type name.
func (p *minicloudProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "minicloud"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
//
//	func (p *minicloudProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
//		resp.Schema = schema.Schema{}
//	}
//
// Schema defines the provider-level schema for configuration data.
func (p *minicloudProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Optional: true,
			},
			"username": schema.StringAttribute{
				Optional: true,
			},
			"password": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
		},
	}
}

// Configure prepares a minicloud API client for data sources and resources.
// func (p *minicloudProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
// }

func (p *minicloudProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {

	tflog.Info(ctx, "Configuring Minicloud client")

	// Retrieve provider data from configuration
	var config minicloudProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown Minicloud API Host",
			"The provider cannot create the Minicloud API client as there is an unknown configuration value for the Minicloud API host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the Minicloud_HOST environment variable.",
		)
	}

	if config.Username.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Unknown Minicloud API Username",
			"The provider cannot create the Minicloud API client as there is an unknown configuration value for the Minicloud API username. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the Minicloud_USERNAME environment variable.",
		)
	}

	if config.Password.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Unknown Minicloud API Password",
			"The provider cannot create the Minicloud API client as there is an unknown configuration value for the Minicloud API password. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the Minicloud_PASSWORD environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	host := os.Getenv("MINICLOUD_HOST")
	username := os.Getenv("MINICLOUD_USERNAME")
	password := os.Getenv("MINICLOUD_PASSWORD")

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	if !config.Username.IsNull() {
		username = config.Username.ValueString()
	}

	if !config.Password.IsNull() {
		password = config.Password.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing Minicloud API Host",
			"The provider cannot create the Minicloud API client as there is a missing or empty value for the Minicloud API host. "+
				"Set the host value in the configuration or use the Minicloud_HOST environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if username == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Missing Minicloud API Username",
			"The provider cannot create the Minicloud API client as there is a missing or empty value for the Minicloud API username. "+
				"Set the username value in the configuration or use the Minicloud_USERNAME environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if password == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Missing Minicloud API Password",
			"The provider cannot create the Minicloud API client as there is a missing or empty value for the Minicloud API password. "+
				"Set the password value in the configuration or use the Minicloud_PASSWORD environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "minicloud_host", host)
	ctx = tflog.SetField(ctx, "minicloud_username", username)
	ctx = tflog.SetField(ctx, "minicloud_password", password)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "minicloud_password") // Mask/hide the password in the log

	tflog.Debug(ctx, "Creating Minicloud client")

	// Create a new minicloud client using the configuration values
	//client, err := minicloud.NewClient(&host, &username, &password)
	client, err := minicloud.NewClient(host, username, password)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Minicloud API Client",
			"An unexpected error occurred when creating the Minicloud API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Minicloud Client Error: "+err.Error(),
		)
		return
	}

	// Make the Minicloud client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured Minicloud client", map[string]any{"success": true})
}

// DataSources defines the data sources implemented in the provider.
//
//	func (p *minicloudProvider) DataSources(_ context.Context) []func() datasource.DataSource {
//		return nil
//	}
//
// DataSources defines the data sources implemented in the provider.
func (p *minicloudProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewVmsDataSource,
	}
}

// Resources defines the resources implemented in the provider.
//
//	func (p *minicloudProvider) Resources(_ context.Context) []func() resource.Resource {
//		return nil
//	}
func (p *minicloudProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewVmsResource,
	}
}
