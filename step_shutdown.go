package main

import (
	"github.com/mitchellh/multistep"
	"github.com/hashicorp/packer/packer"
	"github.com/vmware/govmomi/object"
	"context"
)

type StepShutdown struct{}

func (s *StepShutdown) Run(state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(packer.Ui)
	ui.Say("VM shutdown...")

	vm := state.Get("vm").(*object.VirtualMachine)
	ctx := state.Get("ctx").(context.Context)
	task, err := vm.PowerOff(ctx)
	if err != nil {
		return multistep.ActionHalt
	}
	_, err = task.WaitForResult(ctx, nil)
	if err != nil {
		return multistep.ActionHalt
	}

	ui.Say("VM stopped")
	return multistep.ActionContinue
}

func (s *StepShutdown) Cleanup(multistep.StateBag) {}

