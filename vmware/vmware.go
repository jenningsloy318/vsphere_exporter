package vmware

import (
	"context"
	"fmt"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/soap"
	"net/url"
)

type VMClient struct {
	ctx           context.Context
	govmomiClient *govmomi.Client
}

func NewVMClient(vcURL string, username string, password string, insecure bool) *VMClient {

	ctx := context.Background()
	vsURL, err := soap.ParseURL(vcURL)

	if err != nil {
		fmt.Errorf("error when parsing the vcenter URL,%v", err)
	}

	vsURL.User = url.UserPassword(username, password)

	newVcClient, err := govmomi.NewClient(ctx, vsURL, insecure)

	if err != nil {
		fmt.Errorf("error when creating new vc client ,%v", err)
	}

	err = newVcClient.Login(ctx, vsURL.User)
	if err != nil {
		fmt.Errorf("error when login vc  ,%v", err)
	}
	return &VMClient{
		ctx:           ctx,
		govmomiClient: newVcClient,
	}
}

func (vmc *VMClient) ListVirtualMachine() error {
	vim25Client := vmc.govmomiClient.Client
	ctx := vmc.ctx
	viewManager := view.NewManager(vim25Client)

	virtualMachineListView, err := viewManager.CreateContainerView(ctx, vim25Client.ServiceContent.RootFolder, []string{"VirtualMachine"}, true)

	if err != nil {
		return err
	}

	var virtualMachineList []mo.VirtualMachine
	err = virtualMachineListView.Retrieve(ctx, []string{"VirtualMachine"}, []string{"summary"}, &virtualMachineList)
	if err != nil {
		return err
	}

	for _, vm := range virtualMachineList {
		fmt.Printf("%s: %s\n", vm.Summary.Config.Name, vm.Summary.Config.GuestFullName)
	}
	return nil
}

func (vmc *VMClient) ListHost() error {
	vim25Client := vmc.govmomiClient.Client
	ctx := vmc.ctx
	viewManager := view.NewManager(vim25Client)

	hostSystemListView, err := viewManager.CreateContainerView(ctx, vim25Client.ServiceContent.RootFolder, []string{"HostSystem"}, true)

	if err != nil {
		return err
	}

	var hostSystemList []mo.HostSystem
	err = hostSystemListView.Retrieve(ctx, []string{"HostSystem"}, []string{"summary"}, &hostSystemList)
	if err != nil {
		return err
	}

	for _, hostSystem := range hostSystemList {

		totalCPU := int64(hostSystem.Summary.Hardware.CpuMhz) * int64(hostSystem.Summary.Hardware.NumCpuCores) / 1000
		freeCPU := int64(totalCPU) - int64(hostSystem.Summary.QuickStats.OverallCpuUsage)/1000
		totalMemory := float64(hostSystem.Summary.Hardware.MemorySize) / 1024 / 1024 / 1024
		usedMemory := float64(hostSystem.Summary.QuickStats.OverallMemoryUsage) / 1024
		fmt.Printf("Host %s\t", hostSystem.Summary.Config.Name)
		fmt.Printf("Total CPU: %d GHz\t", totalCPU)
		fmt.Printf("Free CPU: %d GHz\t", freeCPU)
		fmt.Printf("Total Memory: %f GB\t", totalMemory)
		fmt.Printf("Used Memory: %f GB\n", usedMemory)
	}
	return nil
}
