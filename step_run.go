package main

import (
	"github.com/mitchellh/multistep"
	"github.com/hashicorp/packer/packer"
	"github.com/vmware/govmomi/object"
	"context"
	"fmt"
)

type StepRun struct{}

func (s *StepRun) Run(state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(packer.Ui)
	vm := state.Get("vm").(*object.VirtualMachine)
	ctx := state.Get("ctx").(context.Context)

	ui.Say("VM power on...")
	task, err := vm.PowerOn(ctx)
	if err != nil {
		return multistep.ActionHalt
	}
	_, err = task.WaitForResult(ctx, nil)
	if err != nil {
		return multistep.ActionHalt
	}

	ui.Say("VM mounting tools...")
	err = vm.MountToolsInstaller(ctx)
	if err != nil {
		return multistep.ActionHalt
	}

	ui.Say("VM waiting for IP...")
	ip, err := vm.WaitForIP(ctx)
	if err != nil {
		return multistep.ActionHalt
	}

	state.Put("ip", ip)
	ui.Say(fmt.Sprintf("VM ip %v", ip))
	return multistep.ActionContinue
}

func (s *StepRun) Cleanup(multistep.StateBag) {}
