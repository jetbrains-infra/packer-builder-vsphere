package iso

import (
	"github.com/hashicorp/packer/packer"
	"github.com/jetbrains-infra/packer-builder-vsphere/driver"
	"github.com/mitchellh/multistep"
	"fmt"
)

type BootConfig struct {
	BootCommand []string `mapstructure:"boot_command"`
}

func (c *BootConfig) Prepare() []error {
	return nil
}

type StepBootCommand struct {
	Config *BootConfig
}

func (s *StepBootCommand) Run(state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(packer.Ui)
	vm := state.Get("vm").(*driver.VirtualMachine)

	ui.Say("Typing boot command...")
	for _, command := range s.Config.BootCommand {
		_, err := vm.PutUsbScanCodes(command)
		if err != nil {
			state.Put("error", fmt.Errorf("error typing a boot command: %v", err))
			return multistep.ActionHalt
		}
	}

	return multistep.ActionContinue
}

func (s *StepBootCommand) Cleanup(state multistep.StateBag) {}
