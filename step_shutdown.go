package main

import (
	"github.com/mitchellh/multistep"
	"github.com/hashicorp/packer/packer"
	"github.com/vmware/govmomi/object"
	"context"
	"fmt"
	"log"
	"time"
	"bytes"
	"github.com/vmware/govmomi/vim25/types"
)

type StepShutdown struct{
	Command string
}

func (s *StepShutdown) Run(state multistep.StateBag) multistep.StepAction {
	// is set during the communicator.StepConnect
	comm := state.Get("communicator").(packer.Communicator)
	ui := state.Get("ui").(packer.Ui)
	vm := state.Get("vm").(*object.VirtualMachine)
	ctx := state.Get("ctx").(context.Context)

	ui.Say("VM shutdown...")

	if s.Command != "" {
		ui.Say("Gracefully halting virtual machine...")
		log.Printf("Executing shutdown command: %s", s.Command)

		var stdout, stderr bytes.Buffer
		cmd := &packer.RemoteCmd{
			Command: s.Command,
			Stdout:  &stdout,
			Stderr:  &stderr,
		}
		if err := comm.Start(cmd); err != nil {
			err := fmt.Errorf("Failed to send shutdown command: %s", err)
			state.Put("error", err)
			ui.Error(err.Error())
			return multistep.ActionHalt
		}

		var power_state types.VirtualMachinePowerState
		var err error
		// TODO: add timeout
		for !cmd.Exited {
			ui.Say("Waiting for remote cmd to finish...")
			time.Sleep(150 * time.Millisecond)
		}
		if cmd.ExitStatus != 0 {
			err := fmt.Errorf("Cmd exit status %v, not 0", cmd.ExitStatus)
			state.Put("error", err)
			ui.Error(err.Error())
			return multistep.ActionHalt
		}
		for power_state, err = vm.PowerState(ctx); err == nil && power_state != types.VirtualMachinePowerStatePoweredOff; {
			ui.Say("Waiting for VM to finally shut down...")
			time.Sleep(150 * time.Millisecond)
		}
		if err != nil {
			state.Put("error", err)
			ui.Error(err.Error())
			return multistep.ActionHalt
		}

	} else {
		ui.Say("Forcibly halting virtual machine...")

		// TODO: possibly add vm.ShutdownGuest()?
		task, err := vm.PowerOff(ctx)
		if err != nil {
			return multistep.ActionHalt
		}
		_, err = task.WaitForResult(ctx, nil)
		if err != nil {
			return multistep.ActionHalt
		}
	}

	ui.Say("VM stopped")
	return multistep.ActionContinue
}

func (s *StepShutdown) Cleanup(multistep.StateBag) {}

