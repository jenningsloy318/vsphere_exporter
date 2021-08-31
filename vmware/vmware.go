package vmware

import (
	"context"
	"fmt"
	"net/url"

	"github.com/prometheus/common/log"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/soap"
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
	err = virtualMachineListView.Retrieve(ctx, []string{"VirtualMachine"}, []string{"summary"}, &virtualMachineList)
	return virtualMachineList, err
}

func (vmc *VMClient) ListHost() ([]mo.HostSystem, error) {
	vim25Client := vmc.govmomiClient.Client
	ctx := vmc.ctx
	viewManager := view.NewManager(vim25Client)

	hostSystemListView, err := viewManager.CreateContainerView(ctx, vim25Client.ServiceContent.RootFolder, []string{"HostSystem"}, true)
	if err != nil {
		return nil, err
	}

	var hostSystemList []mo.HostSystem
	err = hostSystemListView.Retrieve(ctx, []string{"HostSystem"}, []string{"summary"}, &hostSystemList)
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
	err = datastoreListView.Retrieve(ctx, []string{"Datastore"}, []string{"summary"}, &datastoreList)
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
	err = networkListView.Retrieve(ctx, []string{"Network"}, []string{"summary"}, &networkList)
	return networkList, err

}
func (vmc *VMClient) Logout() error {

	err := vmc.govmomiClient.Logout(vmc.ctx)
	return err

}
