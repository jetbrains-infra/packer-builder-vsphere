package common

import (
	"context"
	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
	"github.com/jetbrains-infra/packer-builder-vsphere/driver"
)

type RunConfig struct {
	BootOrder string `mapstructure:"boot_order"` // example: "floppy,cdrom,ethernet,disk"
}

type StepRun struct {
	Config   *RunConfig
}

func (s *StepRun) Run(_ context.Context, state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(packer.Ui)
	vm := state.Get("vm").(*driver.VirtualMachine)

	ui.Say("Power on VM...")
	err := vm.PowerOn()
	if err != nil {
		state.Put("error", err)
		return multistep.ActionHalt
	}

	return multistep.ActionContinue
}

func (s *StepRun) Cleanup(state multistep.StateBag) {
	ui := state.Get("ui").(packer.Ui)
	vm := state.Get("vm").(*driver.VirtualMachine)

	_, cancelled := state.GetOk(multistep.StateCancelled)
	_, halted := state.GetOk(multistep.StateHalted)
	if !cancelled && !halted {
		return
	}

	ui.Say("Power off VM...")

	err := vm.PowerOff()
	if err != nil {
		ui.Error(err.Error())
	}
}
