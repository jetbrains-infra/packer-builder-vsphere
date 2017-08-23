package driver

import "github.com/vmware/govmomi/object"

func (d *Driver) FindVM(name string) (*object.VirtualMachine, error) {
	return d.finder.VirtualMachine(d.Ctx, name)
}
