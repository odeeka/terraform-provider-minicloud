package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/odeeka/minicloud-client-go/minicloud"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &vmsResource{}
	_ resource.ResourceWithConfigure = &vmsResource{}
)

// NewVmsResource is a helper function to simplify the provider implementation.
func NewVmsResource() resource.Resource {
	return &vmsResource{}
}

// vmsResource is the resource implementation.
type vmsResource struct {
	client *minicloud.Client
}

// vmsResourceModel maps vms schema data.
type vmsResourceModel struct {
	ID          types.Int64   `tfsdk:"id"`
	Name        types.String  `tfsdk:"name"`
	Image       types.String  `tfsdk:"image"`
	CPU         types.Float64 `tfsdk:"cpu"`
	Memory      types.Int64   `tfsdk:"memory"`
	Ports       types.List    `tfsdk:"ports"`
	Env         types.Map     `tfsdk:"env"`
	LastUpdated types.String  `tfsdk:"last_updated"` // Need to set for 'Update' mechanism
}

// Metadata returns the resource type name.
func (r *vmsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vms"
}

// Configure adds the provider configured client to the resource.
func (r *vmsResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = client
}

// Schema defines the schema for the resource.
func (r *vmsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed: true,
				// Keep ID from state during plan to avoid showing "known after apply"
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"image": schema.StringAttribute{
				Required: true,
			},
			"cpu": schema.Float64Attribute{
				Required: true,
			},
			"memory": schema.Int64Attribute{
				Required: true,
			},
			"ports": schema.ListAttribute{
				ElementType: types.Int64Type,
				Optional:    true,
				Description: "List of ports to expose on the VM",
			},
			// Can be changed to 'tags'
			"env": schema.MapAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Environment variables for the VM",
			},
			// Need to addd for 'Update' mechanism
			"last_updated": schema.StringAttribute{
				Computed:    true,
				Description: "Timestamp of the last update",
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *vmsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan vmsResourceModel

	// Read the attributes from plan
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var envMap map[string]string
	plan.Env.ElementsAs(ctx, &envMap, false)

	// Call the client to create the VM
	newVM, err := r.client.CreateVM(minicloud.VM{
		Name:        plan.Name.ValueString(),
		Image:       plan.Image.ValueString(),
		CPU:         plan.CPU.ValueFloat64(),
		Memory:      int(plan.Memory.ValueInt64()),
		Ports:       []int{}, // default
		Env:         envMap,
		ContainerID: "", // default
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating VM",
			fmt.Sprintf("Could not create VM: %s", err),
		)
		return
	}

	plan.ID = types.Int64Value(int64(newVM.ID))
	plan.Name = types.StringValue(newVM.Name)
	plan.Image = types.StringValue(newVM.Image)
	plan.CPU = types.Float64Value(newVM.CPU)
	plan.Memory = types.Int64Value(int64(newVM.Memory))
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	fmt.Printf("[DEBUG] Final value before setting state: last_updated=%v (type: %T)\n", plan.LastUpdated.ValueString(), plan.LastUpdated)

	// Update the state based on newly created VM
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *vmsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state vmsResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the all VMs
	// Replace this block with GetVmByID()
	vms, err := r.client.ListVMs()
	if err != nil {
		resp.Diagnostics.AddError("Error reading VM", fmt.Sprintf("Could not list VMs: %s", err))
		return
	}

	var vm *minicloud.VM
	for _, v := range vms {
		if v.ID == state.ID.ValueInt64() {
			vm = &v
			break
		}
	}

	if vm == nil {
		// Ha a VM már nem létezik, töröljük a state-et
		resp.State.RemoveResource(ctx)
		return
	}

	state.ID = types.Int64Value(int64(vm.ID))
	state.Name = types.StringValue(vm.Name)
	state.Image = types.StringValue(vm.Image)
	state.CPU = types.Float64Value(vm.CPU)
	state.Memory = types.Int64Value(int64(vm.Memory))
	state.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Optional parameters
	envMap, _ := types.MapValueFrom(ctx, types.StringType, vm.Env)
	state.Env = envMap

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
// 1) Retrieves value from the plan
// 2) Generates an API request body from the plan values
// 3) Updates the VM resource
// 4) Maps the response body to resource schema attributes
// 5) Sets the LastUpdated attributes
// 6) Sets Terraform's state with the updated VM
func (r *vmsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan vmsResourceModel
	var state vmsResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.ID = state.ID // IMPORTANT FOR UPDATING

	var portList []int
	port_err := plan.Ports.ElementsAs(ctx, &portList, false)
	if port_err != nil {
		resp.Diagnostics.AddError(
			"Port conversion error",
			fmt.Sprintf("Could not update the ports for VM ID %d: %s", plan.ID.ValueInt64(), port_err),
		)
		return
	}

	var envMap map[string]string
	env_err := plan.Env.ElementsAs(ctx, &envMap, false)
	if env_err != nil {
		resp.Diagnostics.AddError(
			"Env map conversion error",
			fmt.Sprintf("Could not update the env map for VM ID %d: %s", plan.ID.ValueInt64(), env_err),
		)
		return
	}

	updatedVM, err := r.client.UpdateVM(plan.ID.ValueInt64(), minicloud.VM{
		ID:          plan.ID.ValueInt64(),
		Name:        plan.Name.ValueString(),
		Image:       plan.Image.ValueString(),
		CPU:         plan.CPU.ValueFloat64(),
		Memory:      int(plan.Memory.ValueInt64()),
		Ports:       portList,
		Env:         envMap,
		ContainerID: "",
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating VM",
			fmt.Sprintf("Could not update VM ID %d: %s", plan.ID.ValueInt64(), err),
		)
		return
	}

	fmt.Printf("[DEBUG] Update returned VM: ID=%d, Name=%s, CPU=%f, Memory=%d\n", updatedVM.ID, updatedVM.Name, updatedVM.CPU, updatedVM.Memory)

	// Update state with updated values
	plan.ID = types.Int64Value(int64(updatedVM.ID))
	plan.Name = types.StringValue(updatedVM.Name)
	plan.Image = types.StringValue(updatedVM.Image)
	plan.CPU = types.Float64Value(updatedVM.CPU)
	plan.Memory = types.Int64Value(int64(updatedVM.Memory))
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Update the state based on updated VM
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *vmsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state vmsResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing resource (VM)
	err := r.client.DeleteVM(state.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleteing the VM",
			fmt.Sprintf("Could not delete VM ID %d: %s", state.ID.ValueInt64(), err),
		)
		return
	}

}
