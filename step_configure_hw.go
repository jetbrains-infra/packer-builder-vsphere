package main

import (
	"github.com/mitchellh/multistep"
	"github.com/hashicorp/packer/packer"
	"strconv"
	"github.com/vmware/govmomi/vim25/types"
	"context"
	"github.com/vmware/govmomi/object"
)

type StepConfigureHW struct{
	config *Config
}

func (s *StepConfigureHW) Run(state multistep.StateBag) multistep.StepAction {
	vm := state.Get("vm").(*object.VirtualMachine)
	ctx := state.Get("ctx").(context.Context)

	var confSpec types.VirtualMachineConfigSpec
	confNotEmpty := false
	// configure HW
	if s.config.Cpus != "" {
		cpus, err := strconv.Atoi(s.config.Cpus)
		if err != nil {
			state.Put("error", err)
			return multistep.ActionHalt
		}

		confNotEmpty = true
		confSpec.NumCPUs = int32(cpus)
	}
	if s.config.Ram != "" {
		ram, err := strconv.Atoi(s.config.Ram)
		if err != nil {
			state.Put("error", err)
			return multistep.ActionHalt
		}

		confNotEmpty = true
		confSpec.MemoryMB = int64(ram)
	}

	ui := state.Get("ui").(packer.Ui)
	if confNotEmpty {
		ui.Say("configuring virtual hardware...")
		task, err := vm.Reconfigure(ctx, confSpec)
		if err != nil {
			state.Put("error", err)
			return multistep.ActionHalt
		}
		_, err = task.WaitForResult(ctx, nil)
		if err != nil {
			state.Put("error", err)
			return multistep.ActionHalt
		}
	} else {
		ui.Say("skipping the virtual hardware configration...")
	}

	return multistep.ActionContinue
}

func (s *StepConfigureHW) Cleanup(multistep.StateBag) {}
