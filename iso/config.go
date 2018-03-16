package iso

import (
	packerCommon "github.com/hashicorp/packer/common"
	"github.com/hashicorp/packer/helper/communicator"
	"github.com/hashicorp/packer/packer"
	"github.com/hashicorp/packer/template/interpolate"
	"github.com/jetbrains-infra/packer-builder-vsphere/common"
	"fmt"
	"regexp"
)

type Config struct {
	packerCommon.PackerConfig             `mapstructure:",squash"`
	common.RunConfig                      `mapstructure:",squash"`
	BootConfig                            `mapstructure:",squash"`
	common.ConnectConfig                  `mapstructure:",squash"`
	Comm              communicator.Config `mapstructure:",squash"`
	common.ShutdownConfig                 `mapstructure:",squash"`
	CreateSnapshot    bool                `mapstructure:"create_snapshot"`
	ConvertToTemplate bool                `mapstructure:"convert_to_template"`

	CreateConfig `mapstructure:",squash"`
	CDRomConfig  `mapstructure:",squash"`
	FloppyConfig `mapstructure:",squash"`
	ConfigParamsConfig `mapstructure:",squash"`

	ctx interpolate.Context
}

func NewConfig(raws ...interface{}) (*Config, []string, error) {
	c := new(Config)
	if err := common.DecodeConfig(c, &c.ctx, raws...); err != nil {
		return nil, nil, err
	}

	errs := new(packer.MultiError)
	errs = packer.MultiErrorAppend(errs, c.Comm.Prepare(&c.ctx)...)
	errs = packer.MultiErrorAppend(errs, c.RunConfig.Prepare()...)
	errs = packer.MultiErrorAppend(errs, c.ConnectConfig.Prepare()...)
	errs = packer.MultiErrorAppend(errs, c.HardwareConfig.Prepare()...)
	if c.DiskSize <= 0 {
		errs = packer.MultiErrorAppend(errs, fmt.Errorf("'disk_size' must be provided"))
	}
	errs = packer.MultiErrorAppend(errs, c.ShutdownConfig.Prepare()...)
	errs = packer.MultiErrorAppend(errs, c.CreateConfig.Prepare()...)
	if matched, _ := regexp.MatchString(`^vmx-\d+$`, c.Version); !matched && c.Version != "" {
		errs = packer.MultiErrorAppend(errs, fmt.Errorf("'vm_version': " + c.Version + " is not a supported VM version. Valid values start with 'vmx-' and end in a number. Example vmx-10"))
	}

	if len(errs.Errors) > 0 {
		return nil, nil, errs
	}

	return c, nil, nil
}
