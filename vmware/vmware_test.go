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
	newVC, err := NewVMClient(ctx, vsHost, user, pass)
	if err != nil {
		t.Logf("Error when creating vc client, %v", err)
		return
	}
	if newVC.govmomiClient.IsVC() {
		t.Logf("This is a vCenter")
	}
}

func TestVcHost(t *testing.T) {
	ctx := context.Background()
	newVC, err := NewVMClient(ctx, vsHost, user, pass)
	if err != nil {
		t.Logf("Error when creating vc client, %v", err)
		return
	}
	hosts, err := newVC.ListHost()
	if err != nil {
		t.Logf("Error when listing hosts, %v", err)
		return
	}

	t.Logf("hosts %#v\n", hosts[0])
}
func TestVcVM(t *testing.T) {
	ctx := context.Background()
	newVC, err := NewVMClient(ctx, vsHost, user, pass)
	if err != nil {
		t.Logf("Error when creating vc client, %v", err)
		return
	}
	VMs, err := newVC.ListVirtualMachine()
	if err != nil {
		t.Logf("Error when listing virtual machines, %v", err)
		return
	}
	t.Logf("Virtual machines %#v\n", VMs[0])

}

func TestVcDatastore(t *testing.T) {
	ctx := context.Background()
	newVC, err := NewVMClient(ctx, vsHost, user, pass)
	if err != nil {
		t.Logf("Error when creating vc client, %v", err)
		return
	}
	datastores, err := newVC.ListDatastore()
	if err != nil {
		t.Logf("Error when listing datastores, %v", err)
		return
	}

	t.Logf("Datastore %#v\n", datastores[0])

}

func TestVcNetwork(t *testing.T) {
	ctx := context.Background()
	newVC, err := NewVMClient(ctx, vsHost, user, pass)
	if err != nil {
		t.Logf("Error when creating vc client, %v", err)
		return
	}
	networks, err := newVC.ListNetwork()
	if err != nil {
		t.Logf("Error when listing networks, %v", err)
		return
	}

	t.Logf("Networks %#v\n", networks[0])

}
