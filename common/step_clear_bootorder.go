package common

import (
	"context"
	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
	"github.com/jetbrains-infra/packer-builder-vsphere/driver"
)

type StepClearBootOrder struct {
	BootOrder string
	SetOrder  bool
}

func (s *StepClearBootOrder) Run(_ context.Context, state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(packer.Ui)
	vm := state.Get("vm").(*driver.VirtualMachine)

	if s.BootOrder == "" && s.SetOrder {
		ui.Say("Clear boot order...")
		if err := vm.SetBootOrder([]string{"-"}); err != nil {
			state.Put("error", err)
			return multistep.ActionHalt
		}
	}

	return multistep.ActionContinue
}

func (s *StepClearBootOrder) Cleanup(state multistep.StateBag) {}
