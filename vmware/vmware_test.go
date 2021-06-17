package vmware

import (
	"context"
	"testing"
)

var vsHost = "10.36.51.11"
var user = "user"
var pass = "pass"

func TestVC(t *testing.T) {

	ctx := context.Background()
	newVC := NewVMClient(ctx, vsHost, user, pass)
	if newVC.govmomiClient.IsVC() {
		t.Logf("This is a vCenter")
	}
}

func TestVcHost(t *testing.T) {
	ctx := context.Background()
	newVC := NewVMClient(ctx, vsHost, user, pass)

	hosts, err := newVC.ListHost()
	if err != nil {
		return
	}

	t.Logf("hosts %v\n", hosts)
}
func TestVcVM(t *testing.T) {
	ctx := context.Background()
	newVC := NewVMClient(ctx, vsHost, user, pass)

	VMs, err := newVC.ListVirtualMachine()
	if err != nil {
		return
	}
	t.Logf("Virtual machines %v\n", VMs)

}

func TestVcDatastore(t *testing.T) {
	ctx := context.Background()
	newVC := NewVMClient(ctx, vsHost, user, pass)
	datastores, err := newVC.ListDatastore()
	if err != nil {
		return
	}
	t.Logf("Datastore %v\n", datastores)

}

func TestVcNetwork(t *testing.T) {
	ctx := context.Background()
	newVC := NewVMClient(ctx, vsHost, user, pass)

	networks, err := newVC.ListNetwork()
	if err != nil {
		return
	}
	t.Logf("Networks %v\n", networks)

}
