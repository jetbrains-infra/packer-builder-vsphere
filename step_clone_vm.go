package main

import (
	"github.com/vmware/govmomi"
	"context"
	"github.com/mitchellh/multistep"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/govmomi/object"
	"github.com/hashicorp/packer/packer"
)

type StepCloneVM struct{}

func (s *StepCloneVM) Run(state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(packer.Ui)
	ui.Say("start cloning...")

	err := CloneVM(state)
	if err != nil {
		state.Put("error", err)
		return multistep.ActionHalt
	} else {
		ui.Say("finished cloning.")
	}

	return multistep.ActionContinue
}

func (s *StepCloneVM) Cleanup(multistep.StateBag) {}

func CloneVM(state multistep.StateBag) error {
	config := state.Get("config").(*Config)
	vm_params := config.vm_params
	vm_custom := config.vm_custom

	// Prepare entities: client (authentification), finder, folder, virtual machine
	client := state.Get("client").(*govmomi.Client)
	folder := state.Get("folder").(*object.Folder)
	vm_src := state.Get("vm_src").(*object.VirtualMachine)
	ctx := state.Get("ctx").(context.Context)

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

	info, err := task.WaitForResult(ctx, nil)
	if err != nil {
		return err
	}

	vm := object.NewVirtualMachine(client.Client, info.Result.(types.ManagedObjectReference))
	task, err = vm.PowerOn(ctx)
	if err != nil {
		return err
	}
	_, err = task.WaitForResult(ctx, nil)
	if err != nil {
		return err
	}
	err = vm.MountToolsInstaller(ctx)
	if err != nil {
		return err
	}

	result, err := vm.WaitForIP(ctx)
	if err != nil {
		return err
	} else {
		state.Get("ui").(packer.Ui).Say(result)
	}

	return nil
}
