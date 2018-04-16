package iso

import (
	"github.com/hashicorp/packer/packer"
	"github.com/jetbrains-infra/packer-builder-vsphere/driver"
	"github.com/mitchellh/multistep"
	"fmt"
	"github.com/vmware/govmomi/vim25/types"
)

type StepRemoveCDRom struct {
	Config *CDRomConfig
}

func (s *StepRemoveCDRom) Run(state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(packer.Ui)
	vm := state.Get("vm").(*driver.VirtualMachine)

	ui.Say("Removing CDRoms...")
	devices, err := vm.Devices()
	if err != nil {
		state.Put("error", fmt.Errorf("error removing cdrom: error listing devices: %v", err))
		return multistep.ActionHalt
	}
	cdroms := devices.SelectByType((*types.VirtualCdrom)(nil))
	if err = vm.RemoveDevice(false, cdroms...); err != nil {
		state.Put("error", fmt.Errorf("error removing cdrom: %v", err))

		return multistep.ActionHalt
	}

	return multistep.ActionContinue
}

func (s *StepRemoveCDRom) Cleanup(state multistep.StateBag) {
}
