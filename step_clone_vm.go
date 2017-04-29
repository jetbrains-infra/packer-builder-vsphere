package main

import (
	"github.com/vmware/govmomi"
	"context"
	"github.com/mitchellh/multistep"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/govmomi/object"
	"github.com/hashicorp/packer/packer"
)

type CloningEnv struct {
	client *govmomi.Client
	folder *object.Folder
	vm_src *object.VirtualMachine
	ctx    context.Context
}

func NewCloningEnv(state multistep.StateBag) *CloningEnv {
	env := new(CloningEnv)
	env.client = state.Get("client").(*govmomi.Client)
	env.folder = state.Get("folder").(*object.Folder)
	env.vm_src = state.Get("vm_src").(*object.VirtualMachine)
	env.ctx = state.Get("ctx").(context.Context)
	return env
}

type StepCloneVM struct{
	vm_params VMParams
	vm_custom VMCustomParams
}

func (s *StepCloneVM) Run(state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(packer.Ui)
	ui.Say("start cloning...")

	env := NewCloningEnv(state)
	vm, err := CloneVM(s.vm_params, s.vm_custom, env)
	if err != nil {
		state.Put("error", err)
		return multistep.ActionHalt
	}

	state.Put("vm", vm)
	return multistep.ActionContinue
}

func (s *StepCloneVM) Cleanup(multistep.StateBag) {}

func CloneVM(vm_params VMParams, vm_custom VMCustomParams, env *CloningEnv) (vm *object.VirtualMachine, err error) {
	vm = nil
	err = nil

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
	task, err := env.vm_src.Clone(env.ctx, env.folder, vm_params.Vm_target_name, cloneSpec)
	if err != nil {
		return
	}

	info, err := task.WaitForResult(env.ctx, nil)
	if err != nil {
		return
	}

	vm = object.NewVirtualMachine(env.client.Client, info.Result.(types.ManagedObjectReference))
	return vm, nil
}
