package main

import (
	"fmt"

	"github.com/hashicorp/packer/common"
	"github.com/hashicorp/packer/helper/communicator"
	"github.com/hashicorp/packer/helper/config"
	"github.com/hashicorp/packer/packer"
	"github.com/hashicorp/packer/template/interpolate"
	"strconv"
)

type Config struct {
	common.PackerConfig `mapstructure:",squash"`
	communicator.Config `mapstructure:",squash"`

	Url              string `mapstructure:"url"`
	Username         string `mapstructure:"username"`
	Password         string `mapstructure:"password"`

	Template         string `mapstructure:"template"`
	Vm_name          string `mapstructure:"vm_name"`
	Folder_name      string `mapstructure:"folder_name"`
	Dc_name          string `mapstructure:"dc_name"`

	Cpus             string `mapstructure:"cpus"`
	Shutdown_command string `mapstructure:"shutdown_command"`
	Ram              string `mapstructure:"RAM"`
	//TODO: add more options

	ctx      interpolate.Context
}

func NewConfig(raws ...interface{}) (*Config, []string, error) {
	c := new(Config)
	err := config.Decode(c, &config.DecodeOpts{
		Interpolate:        true,
		InterpolateContext: &c.ctx,
	}, raws...)
	if err != nil {
		return nil, nil, err
	}

	// Accumulate any errors
	errs := new(packer.MultiError)

	// Prepare config(s)
	errs = packer.MultiErrorAppend(errs, c.Config.Prepare(&c.ctx)...)

	// Check the required params
	if c.Url == "" {
		errs = packer.MultiErrorAppend(errs, fmt.Errorf("URL required"))
	}
	if c.Username == "" {
		errs = packer.MultiErrorAppend(errs, fmt.Errorf("Username required"))
	}
	if c.Password == "" {
		errs = packer.MultiErrorAppend(errs, fmt.Errorf("Password required"))
	}
	if c.Template == "" {
		errs = packer.MultiErrorAppend(errs, fmt.Errorf("Template VM name required"))
	}
	if c.Vm_name == "" {
		errs = packer.MultiErrorAppend(errs, fmt.Errorf("Target VM name required"))
	}

	// Verify numeric parameters if present
	if c.Cpus != "" {
		if _, err = strconv.Atoi(c.Cpus); err != nil {
			errs = packer.MultiErrorAppend(errs, fmt.Errorf("Invalid number of cpu sockets"))
		}
	}
	if c.Ram != "" {
		if _, err = strconv.Atoi(c.Ram); err != nil {
			errs = packer.MultiErrorAppend(errs, fmt.Errorf("Invalid number for Ram"))
		}
	}

	// Warnings
	var warnings []string
	if c.Shutdown_command == "" {
		warnings = append(warnings,
			"A shutdown_command was not specified. Without a shutdown command, Packer\n"+
				"will forcibly halt the virtual machine, which may result in data loss.")
	}

	if len(errs.Errors) > 0 {
		return nil, warnings, errs
	}

	return c, warnings, nil
}
