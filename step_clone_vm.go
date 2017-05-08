package main

import (
	"github.com/vmware/govmomi"
	"context"
	"github.com/mitchellh/multistep"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/govmomi/object"
	"github.com/hashicorp/packer/packer"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/vim25/mo"
	"fmt"
	"net/url"
)

type CloneParameters struct {
	client   *govmomi.Client
	folder   *object.Folder
	vm_src   *object.VirtualMachine
	ctx       context.Context
	config   *Config
	confSpec *types.VirtualMachineConfigSpec
}

type StepCloneVM struct{
	config *Config
	success bool
}

func (s *StepCloneVM) Run(state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(packer.Ui)
	ui.Say("start cloning...")

	confSpec := state.Get("confSpec").(types.VirtualMachineConfigSpec)

	// Prepare entities: client (authentification), finder, folder, virtual machine
	client, ctx, err := createClient(s.config.Url, s.config.Username, s.config.Password)
	if err != nil {
		state.Put("error", err)
		return multistep.ActionHalt
	}
	finder, ctx, err := createFinder(ctx, client, s.config.Dc_name)
	if err != nil {
		state.Put("error", err)
		return multistep.ActionHalt
	}
	folder, err := finder.FolderOrDefault(ctx, s.config.Folder_name)
	if err != nil {
		state.Put("error", err)
		return multistep.ActionHalt
	}
	vm_src, err := finder.VirtualMachine(ctx, s.config.Template)
	if err != nil {
		state.Put("error", err)
		return multistep.ActionHalt
	}

	cloneParameters := CloneParameters{
		client: client,
		folder: folder,
		vm_src: vm_src,
		ctx: ctx,
		config: s.config,
		confSpec: &confSpec,
	}

	vm, err := cloneVM(&cloneParameters)
	if err != nil {
		state.Put("error", err)
		return multistep.ActionHalt
	}

	state.Put("vm", vm)
	state.Put("ctx", ctx)
	s.success = true
	return multistep.ActionContinue
}

func (s *StepCloneVM) Cleanup(state multistep.StateBag) {
	if !s.success {
		return
	}

	_, cancelled := state.GetOk(multistep.StateCancelled)
	_, halted := state.GetOk(multistep.StateHalted)

	if cancelled || halted {
		vm := state.Get("vm").(*object.VirtualMachine)
		ctx := state.Get("ctx").(context.Context)
		ui := state.Get("ui").(packer.Ui)

		ui.Say("destroying VM...")

		task, err := vm.Destroy(ctx)
		if err != nil {
			ui.Error(err.Error())
			return
		}
		_, err = task.WaitForResult(ctx, nil)
		if err != nil {
			ui.Error(err.Error())
			return
		}
	}
}

func cloneVM(params *CloneParameters) (vm *object.VirtualMachine, err error) {
	vm = nil
	err = nil

	// Creating specs for cloning
	var relocateSpec types.VirtualMachineRelocateSpec
	cloneSpec := types.VirtualMachineCloneSpec{
		Location: relocateSpec,
		Config:   params.confSpec,
		PowerOn:  false,
	}

	// Cloning itself
	task, err := params.vm_src.Clone(params.ctx, params.folder, params.config.Vm_name, cloneSpec)
	if err != nil {
		return
	}

	info, err := task.WaitForResult(params.ctx, nil)
	if err != nil {
		return
	}

	vm = object.NewVirtualMachine(params.client.Client, info.Result.(types.ManagedObjectReference))
	return vm, nil
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
		finder.SetDatacenter(dc)
	} else {
		dc, err := finder.Datacenter(ctx, dc_name)
		if err != nil {
			return nil, nil, err
		}
		finder.SetDatacenter(dc)
	}
	return finder, ctx, nil
}
