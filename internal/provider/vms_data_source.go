package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/odeeka/minicloud-client-go/minicloud"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &vmsDataSource{}
	_ datasource.DataSourceWithConfigure = &vmsDataSource{}
)

// NewvmsDataSource is a helper function to simplify the provider implementation.
func NewVmsDataSource() datasource.DataSource {
	return &vmsDataSource{}
}

// vmsDataSource is the data source implementation.
type vmsDataSource struct {
	client *minicloud.Client
}

// vmsDataSourceModel maps the data source schema data.
type vmsDataSourceModel struct {
	Vms []vmsModel `tfsdk:"vms"`
}

// vmsModel maps coffees schema data.
type vmsModel struct {
	ID     types.Int64   `tfsdk:"id"`
	Name   types.String  `tfsdk:"name"`
	Image  types.String  `tfsdk:"image"`
	CPU    types.Float64 `tfsdk:"cpu"`
	Memory types.Int64   `tfsdk:"memory"`
}

// Metadata returns the data source type name.
func (d *vmsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vms"
}

// Schema defines the schema for the data source.
//
//	func (d *vmsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
//		resp.Schema = schema.Schema{}
//	}
//
// Schema defines the schema for the data source.
func (d *vmsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"vms": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"image": schema.StringAttribute{
							Computed: true,
						},
						"cpu": schema.Float64Attribute{
							Computed: true,
						},
						"memory": schema.Int64Attribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
// func (d *vmsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
// }
// Read refreshes the Terraform state with the latest data.
func (d *vmsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state vmsDataSourceModel

	vms, err := d.client.ListVMs()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Minicloud VMs",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, vm := range vms {
		vmState := vmsModel{
			ID:     types.Int64Value(int64(vm.ID)),
			Name:   types.StringValue(vm.Name),
			Image:  types.StringValue(vm.Image),
			CPU:    types.Float64Value(vm.CPU),
			Memory: types.Int64Value(int64(vm.Memory)),
		}

		state.Vms = append(state.Vms, vmState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *vmsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*minicloud.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *minicloud.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}
