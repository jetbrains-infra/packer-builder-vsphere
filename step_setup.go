package main

import (
	"github.com/mitchellh/multistep"
	"github.com/hashicorp/packer/packer"
	"github.com/vmware/govmomi/find"
)

type StepSetup struct{
	config *Config
}

func (s *StepSetup) Run(state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(packer.Ui)
	ui.Say("setup...")

	// Prepare entities: client (authentification), finder, folder, virtual machine
	client, ctx, err := createClient(s.config.Url, s.config.Username, s.config.Password)
	if err != nil {
		state.Put("error", err)
		return multistep.ActionHalt
	}

	// Set up finder
	finder := find.NewFinder(client.Client, false)
	dc, err := finder.DatacenterOrDefault(ctx, s.config.DCName)
	if err != nil {
		state.Put("error", err)
		return multistep.ActionHalt
	}
	finder.SetDatacenter(dc)

	// Get source VM
	vmSrc, err := finder.VirtualMachine(ctx, s.config.Template)
	if err != nil {
		state.Put("error", err)
		return multistep.ActionHalt
	}

	state.Put("client", client)
	state.Put("ctx", ctx)
	state.Put("finder", finder)
	state.Put("dc", dc)
	state.Put("vmSrc", vmSrc)
	return multistep.ActionContinue
}

func (s *StepSetup) Cleanup(state multistep.StateBag) {}
