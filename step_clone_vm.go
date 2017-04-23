package main

import (
	"github.com/vmware/govmomi"
	"context"
	"net/url"
	"fmt"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/mitchellh/multistep"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/govmomi/object"
)

type StepCloneVM struct{}

func (s *StepCloneVM) Run(state multistep.StateBag) multistep.StepAction {
	config := state.Get("config").(*Config)
	vm_params := config.vm_params
	vm_custom := config.vm_custom
	err := CloneVM(vm_params, vm_custom)
	if err != nil {
		state.Put("error", err)
		return multistep.ActionHalt
	}
	return multistep.ActionContinue
}

func (s *StepCloneVM) Cleanup(multistep.StateBag) {}

func CloneVM(vm_params VMParams, vm_custom VMCustomParams) error {
	// Prepare entities: client (authentification), finder, folder, virtual machine
	client, ctx, err := createClient(vm_params.Url, vm_params.Username, vm_params.Password)
	if err != nil {
		return err
	}
	finder, ctx, err := createFinder(ctx, client, vm_params.Dc_name)
	if err != nil {
		return err
	}
	folder, err := finder.FolderOrDefault(ctx, vm_params.Folder_name)
	if err != nil {
		return err
	}
	vm_src, ctx, err := findVM_by_name(ctx, finder, vm_params.Vm_source_name)
	if err != nil {
		return err
	}

	// Creating spec's for cloning
	var relocateSpec types.VirtualMachineRelocateSpec

	var confSpec types.VirtualMachineConfigSpec
	// configure HW
	if vm_custom.Cpu_sockets != Unspecified {
		confSpec.NumCPUs = int32(vm_custom.Cpu_sockets)
	}
	if vm_custom.Ram != Unspecified {
		confSpec.MemoryMB = int64(vm_custom.Ram)
	}

	cloneSpec := types.VirtualMachineCloneSpec{
		Location: relocateSpec,
		Config:   &confSpec,
		PowerOn:  false,
	}

	// Cloning itself
	task, err := vm_src.Clone(ctx, folder, vm_params.Vm_target_name, cloneSpec)
	if err != nil {
		return err
	}
	_, err = task.WaitForResult(ctx, nil)
	if err != nil {
		return err
	}

	return nil
}
func createClient(URL, username, password string) (*govmomi.Client, context.Context, error) {
	// create context
	ctx := context.TODO() // an empty, default context (for those, who is unsure)

	// create a client
	// (connected to the specified URL,
	// logged in with the username-password)
	u, err := url.Parse(URL) // create a URL object from string
	if err != nil {
		return nil, nil, err
	}
	u.User = url.UserPassword(username, password) // set username and password for automatical authentification
	fmt.Println(u.String())
	client, err := govmomi.NewClient(ctx, u,true) // creating a client (logs in with given uname&pswd)
	if err != nil {
		return nil, nil, err
	}

	return client, ctx, nil
}

func createFinder(ctx context.Context, client *govmomi.Client, dc_name string) (*find.Finder, context.Context, error) {
	// Create a finder to search for a vm with the specified name
	finder := find.NewFinder(client.Client, false)
	// Need to specify the datacenter
	if dc_name == "" {
		dc, err := finder.DefaultDatacenter(ctx)
		if err != nil {
			return nil, nil, err
		}
		var dc_mo mo.Datacenter
		err = dc.Properties(ctx, dc.Reference(), []string{"name"}, &dc_mo)
		if err != nil {
			return nil, nil, err
		}
		dc_name = dc_mo.Name
		finder.SetDatacenter(dc)
	} else {
		dc, err := finder.Datacenter(ctx, fmt.Sprintf("/%v", dc_name))
		if err != nil {
			return nil, nil, err
		}
		finder.SetDatacenter(dc)
	}
	return finder, ctx, nil
}

func findVM_by_name(ctx context.Context, finder *find.Finder, vm_name string) (*object.VirtualMachine, context.Context, error) {
	vm_o, err := finder.VirtualMachine(ctx, vm_name)
	if err != nil {
		return nil, nil, err
	}
	return vm_o, ctx, nil
}
