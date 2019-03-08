package common

import (
	"context"
	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
	"github.com/jetbrains-infra/packer-builder-vsphere/driver"
	"strings"
)

type StepSetBootOrder struct {
	BootOrder string
	SetOrder bool
}

func (s *StepSetBootOrder) Run(_ context.Context, state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(packer.Ui)
	vm := state.Get("vm").(*driver.VirtualMachine)

	if s.BootOrder != "" {
		ui.Say("Set boot order...")
		order := strings.Split(s.BootOrder, ",")
		if err := vm.SetBootOrder(order); err != nil {
			state.Put("error", err)
			return multistep.ActionHalt
		}
	} else {
		if s.SetOrder {
			ui.Say("Set boot order temporary...")
			if err := vm.SetBootOrder([]string{"disk", "cdrom"}); err != nil {
				state.Put("error", err)
				return multistep.ActionHalt
			}
		}
	}


	return multistep.ActionContinue
}

func (s *StepSetBootOrder) Cleanup(state multistep.StateBag) {}