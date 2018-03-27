package common

import "fmt"

type VMConfig struct {
	VMName           string `mapstructure:"vm_name"`
	Folder           string `mapstructure:"folder"`
	Cluster          string `mapstructure:"cluster"`
	Host             string `mapstructure:"host"`
	ResourcePool     string `mapstructure:"resource_pool"`
	Datastorecluster string `mapstructure:"datastorecluster"`
	Datastore        string `mapstructure:"datastore"`
}

func (c *VMConfig) Prepare() []error {
	var errs []error

	if c.VMName == "" {
		errs = append(errs, fmt.Errorf("Target VM name is required"))
	}
	if c.Cluster == "" && c.Host == "" {
		errs = append(errs, fmt.Errorf("vSphere host or cluster is required"))
	}
	if c.Datastorecluster != "" && c.Datastore != "" {
		errs = append(errs, fmt.Errorf("Datastorecluster and datastore cannot be set at the same time"))
	}

	return errs
}
