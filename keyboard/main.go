package main

import (
	"github.com/jetbrains-infra/packer-builder-vsphere/driver"
	"fmt"
)

func main() {
	d, err := driver.NewDriver(&driver.ConnectConfig{
		VCenterServer:      "vcenter.vsphere65.test",
		Username:           "root",
		Password:           "jetbrains",
		InsecureConnection: true,
	})
	if err != nil {
		panic(err)
	}

	vm, err := d.FindVM("alpine-1")
	if err != nil {
		panic(err)
	}

	n, err := vm.PutUsbScanCodes("abc123")
	if err != nil {
		panic(err)
	}

	fmt.Println(n)
}
