package main

import (
	"github.com/mitchellh/multistep"
	"github.com/hashicorp/packer/packer"
	"context"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/vim25/mo"
	"fmt"
	"net/url"
	"github.com/vmware/govmomi/object"
)

type StepSetupCloningEnv struct{
	vm_params VMParams
}

func (s *StepSetupCloningEnv) Run(state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(packer.Ui)
	ui.Say("setup cloning environment...")

	// Prepare entities: client (authentification), finder, folder, virtual machine
	client, ctx, err := createClient(s.vm_params.Url, s.vm_params.Username, s.vm_params.Password)
	if err != nil {
		state.Put("error", err)
		return multistep.ActionHalt
	}
	finder, ctx, err := createFinder(ctx, client, s.vm_params.Dc_name)
	if err != nil {
		state.Put("error", err)
		return multistep.ActionHalt
	}
	folder, err := finder.FolderOrDefault(ctx, s.vm_params.Folder_name)
	if err != nil {
		state.Put("error", err)
		return multistep.ActionHalt
	}
	vm_src, ctx, err := findVM_by_name(ctx, finder, s.vm_params.Vm_source_name)
	if err != nil {
		state.Put("error", err)
		return multistep.ActionHalt
	}

	state.Put("client", client)
	state.Put("folder", folder)
	state.Put("vm_src", vm_src)
	state.Put("ctx", ctx)

	return multistep.ActionContinue
}

func (s *StepSetupCloningEnv) Cleanup(multistep.StateBag) {}

// TODO: make a separate step
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

// TODO: make a separate step
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

// TODO: make a separate step
func findVM_by_name(ctx context.Context, finder *find.Finder, vm_name string) (*object.VirtualMachine, context.Context, error) {
	vm_o, err := finder.VirtualMachine(ctx, vm_name)
	if err != nil {
		return nil, nil, err
	}
	return vm_o, ctx, nil
}
