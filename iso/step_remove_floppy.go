package iso

import (
	"github.com/mitchellh/multistep"
	"github.com/jetbrains-infra/packer-builder-vsphere/driver"
	"fmt"
	"github.com/vmware/govmomi/vim25/types"
)

type StepRemoveFloppy struct {
	Config    *FloppyConfig
	Datastore string
	Host      string

	uploadedFloppyPath string
}

func (s *StepRemoveFloppy) Run(state multistep.StateBag) multistep.StepAction {
	err := s.runImpl(state)
	if err != nil {
		state.Put("error", err)
		return multistep.ActionHalt
	}
	return multistep.ActionContinue
}

func (s *StepRemoveFloppy) runImpl(state multistep.StateBag) error {
	vm := state.Get("vm").(*driver.VirtualMachine)
	d := state.Get("driver").(*driver.Driver)

	devices, err := vm.Devices()
	if err != nil {
		return fmt.Errorf("error removing floppy: %v", err)
	}
	cdroms := devices.SelectByType((*types.VirtualFloppy)(nil))
	if err = vm.RemoveDevice(false, cdroms...); err != nil {
		return fmt.Errorf("error removing floppy: %v", err)
	}

	if s.uploadedFloppyPath != "" {
		ds, err := d.FindDatastore(s.Datastore, s.Host)
		if err != nil {
			return fmt.Errorf("error removing floppy: %v", err)
		}
		if err := ds.Delete(s.uploadedFloppyPath); err != nil {
			return fmt.Errorf("error deleting floppy image '%v': %v", s.uploadedFloppyPath, err.Error())
		}
	}
	return nil
}

func (s *StepRemoveFloppy) Cleanup(state multistep.StateBag) {
}
