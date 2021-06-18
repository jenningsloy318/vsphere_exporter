package vmware

import (
	"context"
	"fmt"
	"net/url"

	"github.com/prometheus/common/log"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/performance"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/soap"
	"github.com/vmware/govmomi/vim25/types"
)

type VMClient struct {
	ctx           context.Context
	govmomiClient *govmomi.Client
}

func NewVMClient(context context.Context, vcHost string, username string, password string) (*VMClient, error) {

	vcURL, err := soap.ParseURL(fmt.Sprintf("https://%s", vcHost))

	if err != nil {
		log.Errorf("error when parsing the vCenter URL, %v", err)
		return nil, err
	}

	vcURL.User = url.UserPassword(username, password)

	newVcClient, err := govmomi.NewClient(context, vcURL, true)

	if err != nil {
		log.Errorf("error when creating new vCenter client, %v", err)
		return nil, err
	}
	return &VMClient{
		ctx:           context,
		govmomiClient: newVcClient,
	}, nil
}

func (vmc *VMClient) ListVirtualMachine() ([]mo.VirtualMachine, error) {
	vim25Client := vmc.govmomiClient.Client
	ctx := vmc.ctx
	viewManager := view.NewManager(vim25Client)

	virtualMachineListView, err := viewManager.CreateContainerView(ctx, vim25Client.ServiceContent.RootFolder, []string{"VirtualMachine"}, true)
	if err != nil {
		return nil, err
	}

	var virtualMachineList []mo.VirtualMachine
	//https://code.vmware.com/apis/358/vsphere/doc/vim.VirtualMachine.html, VirtualMachine has multiple properties, but here we choose "summary","config","guest","guestHeartbeatStatus","runtime"
	err = virtualMachineListView.Retrieve(ctx, []string{"VirtualMachine"}, []string{"summary", "config", "guest", "guestHeartbeatStatus", "runtime"}, &virtualMachineList)
	return virtualMachineList, err
}

func (vmc *VMClient) ListHost() ([]mo.HostSystem, error) {
	vim25Client := vmc.govmomiClient.Client
	ctx := vmc.ctx
	viewManager := view.NewManager(vim25Client)
	// containerView has multiple properties when creating the container view, mainy categoried as 5 types,  Folder / Datacenter /ComputeResource /ResourcePool /HostSystem,here we can use RootFolder, which is a Folder, for full list of supported property

	hostSystemListView, err := viewManager.CreateContainerView(ctx, vim25Client.ServiceContent.RootFolder, []string{"HostSystem"}, true)
	if err != nil {
		return nil, err
	}

	var hostSystemList []mo.HostSystem
	// https://code.vmware.com/apis/358/vsphere/doc/vim.HostSystem.html, here HostSystem has multiple Properties that can be retrieved, but here we choose "summary","runtime","hardware","config","capability"
	err = hostSystemListView.Retrieve(ctx, []string{"HostSystem"}, []string{"summary", "runtime", "hardware", "config", "capability"}, &hostSystemList)
	return hostSystemList, err

}

func (vmc *VMClient) ListDatastore() ([]mo.Datastore, error) {
	vim25Client := vmc.govmomiClient.Client
	ctx := vmc.ctx
	viewManager := view.NewManager(vim25Client)

	datastoreListView, err := viewManager.CreateContainerView(ctx, vim25Client.ServiceContent.RootFolder, []string{"Datastore"}, true)
	if err != nil {
		return nil, err
	}

	var datastoreList []mo.Datastore
	//https://code.vmware.com/apis/358/vsphere/doc/vim.Datastore.html, datastore have several properties, we choose "summary","info"
	err = datastoreListView.Retrieve(ctx, []string{"Datastore"}, []string{"summary", "info"}, &datastoreList)
	return datastoreList, err

}

func (vmc *VMClient) ListNetwork() ([]mo.Network, error) {
	vim25Client := vmc.govmomiClient.Client
	ctx := vmc.ctx
	viewManager := view.NewManager(vim25Client)

	networkListView, err := viewManager.CreateContainerView(ctx, vim25Client.ServiceContent.RootFolder, []string{"Network"}, true)
	if err != nil {
		return nil, err
	}

	var networkList []mo.Network
	// https://code.vmware.com/apis/358/vsphere/doc/vim.Network.html, network only have four properties, we choose "name" and "summary"
	err = networkListView.Retrieve(ctx, []string{"Network"}, []string{"summary", "name"}, &networkList)
	return networkList, err

}

func (vmc *VMClient) ListPerfCounters() (map[string]*types.PerfCounterInfo, error) {
	vim25Client := vmc.govmomiClient.Client
	ctx := vmc.ctx

	perfManager := performance.NewManager(vim25Client)
	perfCounters, err := perfManager.CounterInfoByName(ctx)
	if err != nil {
		log.Errorf("error when getting perf counters from vcenter, %v", err)
		return nil, err
	}
	return perfCounters, nil

}

func (vmc *VMClient) Logout() error {

	err := vmc.govmomiClient.Logout(vmc.ctx)
	return err

}
