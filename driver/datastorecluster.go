package driver

import (
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/govmomi/vim25/mo"
)

type Datastorecluster struct {
	dsc    *object.StoragePod
	driver *Driver
}

func (d *Driver) NewDatastorecluster(ref *types.ManagedObjectReference) *Datastorecluster {
	return &Datastorecluster{
		dsc:    object.NewStoragePod(d.client.Client, *ref),
		driver: d,
	}
}

func (d *Driver) FindDatastorecluster(name string) (*Datastorecluster, error) {
	dsc, err := d.finder.DatastoreCluster(d.ctx, name)
	if err != nil {
		return nil, err
	}

	return &Datastorecluster{
		dsc:    dsc,
		driver: d,
	}, nil
}

func (dsc *Datastorecluster) Info(params ...string) (*mo.StoragePod, error) {
	var p []string
	if len(params) == 0 {
		p = []string{"*"}
	} else {
		p = params
	}
	var info mo.StoragePod
	err := dsc.dsc.Properties(dsc.driver.ctx, dsc.dsc.Reference(), p, &info)
	if err != nil {
		return nil, err
	}
	return &info, nil
}

func (dsc *Datastorecluster) Name() string {
	return dsc.dsc.Name()
}