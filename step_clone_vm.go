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
	config *Config
	success bool
}

func (s *StepCloneVM) Run(state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(packer.Ui)
	ui.Say("start cloning...")
	confSpec := state.Get("confSpec").(types.VirtualMachineConfigSpec)

	env := NewCloningEnv(state)
	vm, err := CloneVM(s.config, env, &confSpec)
	if err != nil {
		state.Put("error", err)
		return multistep.ActionHalt
	}

	state.Put("vm", vm)
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

func CloneVM(config *Config, env *CloningEnv, confSpec *types.VirtualMachineConfigSpec) (vm *object.VirtualMachine, err error) {
	vm = nil
	err = nil

	// Creating specs for cloning
	var relocateSpec types.VirtualMachineRelocateSpec
	cloneSpec := types.VirtualMachineCloneSpec{
		Location: relocateSpec,
		Config:   confSpec,
		PowerOn:  false,
	}

	// Cloning itself
	task, err := env.vm_src.Clone(env.ctx, env.folder, config.Vm_name, cloneSpec)
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
