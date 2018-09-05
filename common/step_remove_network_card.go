package common

import (
	"context"
	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
	"github.com/jetbrains-infra/packer-builder-vsphere/driver"
	"github.com/vmware/govmomi/vim25/types"
)

type StepRemoveNetworkCard struct {
	RemoveNetworkCard bool
}

func (s *StepRemoveNetworkCard) Run(_ context.Context, state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(packer.Ui)
	vm := state.Get("vm").(*driver.VirtualMachine)

	if s.RemoveNetworkCard {
		ui.Say("Deleting Network Cards...")
		devices, err := vm.Devices()
		if err != nil {
			state.Put("error", err)
			return multistep.ActionHalt
		}

		cards := devices.SelectByType((*types.VirtualEthernetCard)(nil))

		if err = vm.RemoveDevice(true, cards...); err != nil {
			state.Put("error", err)
			return multistep.ActionHalt
		}
	}

	return multistep.ActionContinue
}

func (s *StepRemoveNetworkCard) Cleanup(state multistep.StateBag) {}
