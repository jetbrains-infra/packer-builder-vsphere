package iso

import (
	"context"
	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
	"github.com/jetbrains-infra/packer-builder-vsphere/driver"
	"github.com/vmware/govmomi/vim25/types"
)

type StepRemoveCDRom struct{}

func (s *StepRemoveCDRom) Run(_ context.Context, state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(packer.Ui)
	vm := state.Get("vm").(*driver.VirtualMachine)

	ui.Say("Deleting CD-ROM drives...")
	devices, err := vm.Devices()
	if err != nil {
		state.Put("error", err)
		return multistep.ActionHalt
	}
	cdroms := devices.SelectByType((*types.VirtualCdrom)(nil))
	if err = vm.RemoveDevice(true, cdroms...); err != nil {
		state.Put("error", err)
		return multistep.ActionHalt
	}

	ui.Say("Deleting SATA controller...")
	sata := devices.SelectByType((*types.VirtualAHCIController)(nil))
	if err = vm.RemoveDevice(true, sata...); err != nil {
		state.Put("error", err)
		return multistep.ActionHalt
	}

	return multistep.ActionContinue
}

func (s *StepRemoveCDRom) Cleanup(state multistep.StateBag) {}
