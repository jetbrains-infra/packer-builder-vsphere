package main

import (
	"fmt"

	vmwcommon "github.com/hashicorp/packer/builder/vmware/common"
	"github.com/hashicorp/packer/common"
	"github.com/hashicorp/packer/helper/communicator"
	"github.com/hashicorp/packer/helper/config"
	"github.com/hashicorp/packer/packer"
	"github.com/hashicorp/packer/template/interpolate"
	"strconv"
)

const Unspecified = -1

type VMCustomParams struct {
	Cpu_sockets int
	Ram         int
	// TODO: add more options
}

type VMParams struct {
	Url            string
	Username       string
	Password       string
	Dc_name        string
	Folder_name    string
	Vm_source_name string
	Vm_target_name string
}

type Config struct {
	common.PackerConfig `mapstructure:",squash"`
	CommConfig communicator.Config `mapstructure:",squash"`
	vmwcommon.DriverConfig   `mapstructure:",squash"`
	vmwcommon.SSHConfig      `mapstructure:",squash"`

	Url            string `mapstructure:"url"`
	Username       string `mapstructure:"username"`
	Password       string `mapstructure:"password"`
	Ssh_username   string `mapstructure:"ssh_username"`
	Ssh_password   string `mapstructure:"ssh_password"`

	Dc_name        string `mapstructure:"dc_name"`
	Vm_source_name string `mapstructure:"template"`
	Vm_target_name string `mapstructure:"vm_name"`

	Cpu_sockets    string `mapstructure:"cpus"`
	Shutdown_cmd   string `mapstructure:"shutdown_command"`
	Ram            string `mapstructure:"RAM"`

	// internal
	vm_params VMParams
	vm_custom VMCustomParams

	ctx      interpolate.Context
}


func NewConfig(raws ...interface{}) (*Config, []string, error) {
	c := new(Config)
	err := config.Decode(c, &config.DecodeOpts{
		Interpolate:        true,
		InterpolateContext: &c.ctx,
		InterpolateFilter: &interpolate.RenderFilter{
			Exclude: []string{
				"shutdown_command",
			},
		},
	}, raws...)
	if err != nil {
		return nil, nil, err
	}

	// Accumulate any errors
	errs := new(packer.MultiError)

	// Prepare config(s)
	errs = packer.MultiErrorAppend(errs, c.CommConfig.Prepare(&c.ctx)...)
	errs = packer.MultiErrorAppend(errs, c.DriverConfig.Prepare(&c.ctx)...)
	errs = packer.MultiErrorAppend(errs, c.SSHConfig.Prepare(&c.ctx)...)

	// Check the required params
	templates := map[string]*string{
		"url":            &c.Url,
		"username":       &c.Username,
		"password":       &c.Password,
		"vm_source_name": &c.Vm_source_name,
	}
	for key, ptr := range templates {
		if *ptr == "" {
			errs = packer.MultiErrorAppend(
				errs, fmt.Errorf("%s must be set, %s is present", key, *ptr))
		}
	}

	// Warnings
	var warnings []string
	if c.Shutdown_cmd == "" {
		warnings = append(warnings,
			"A shutdown_command was not specified. Without a shutdown command, Packer\n"+
				"will forcibly halt the virtual machine, which may result in data loss.")
	}

	if len(errs.Errors) > 0 {
		return nil, warnings, errs
	}

	// Set defaults
	if c.Vm_target_name == "" {
		c.Vm_target_name = c.Vm_source_name + "_clone"
	}

	// Set custom params
	c.vm_custom.Cpu_sockets = Unspecified
	if c.Cpu_sockets != "" {
		c.vm_custom.Cpu_sockets, err = strconv.Atoi(c.Cpu_sockets)
		if err != nil {
			return nil, warnings, err
		}
	}
	c.vm_custom.Ram = Unspecified
	if c.Ram != "" {
		c.vm_custom.Ram, err = strconv.Atoi(c.Ram)
		if err != nil {
			return nil, warnings, err
		}
	}

	// Set required params
	c.vm_params.Url = c.Url
	c.vm_params.Username = c.Username
	c.vm_params.Password = c.Password
	c.vm_params.Dc_name = c.Dc_name
	c.vm_params.Vm_source_name = c.Vm_source_name
	c.vm_params.Vm_target_name = c.Vm_target_name

	return c, warnings, nil
}
