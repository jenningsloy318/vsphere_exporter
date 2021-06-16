package vmware

import (
	"testing"
)

var testurl = ""
var user = ""
var pass = ""

func TestVC(t *testing.T) {

	newVC := NewVMClient(testurl, user, pass, true)

	t.Logf("This is VC: %v", newVC.govmomiClient.IsVC())
}

func TestVcHosts(t *testing.T) {

	newVC := NewVMClient(testurl, user, pass, true)
	newVC.ListHost()
}
func TestVcVMs(t *testing.T) {

	newVC := NewVMClient(testurl, user, pass, true)
	newVC.ListVirtualMachine()
}
