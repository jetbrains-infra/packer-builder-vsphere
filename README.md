# Packer Builder for VMware vSphere

This repo contains two plugins for [HashiCorp Packer](https://www.packer.io/). The plugins use the native vSphere API, and create virtual machines remotely. The `vsphere-iso` plugin
is used to create VMWare templates from scratch. The `vsphere-clone` plugin is used to create VMWare templates from an existing base template.

- VMware Player is not required
- Official vCenter API is used, no ESXi host [modification](https://www.packer.io/docs/builders/vmware-iso.html#building-on-a-remote-vsphere-hypervisor) is required 

## Usage
* Download the plugins from the [releases page](https://github.com/jetbrains-infra/packer-builder-vsphere/releases).
* [Install](https://www.packer.io/docs/extending/plugins.html#installing-plugins) the plugins, or simply put them into the same directory with configuration files. On Linux and macOS run `chmod +x` on the plugin binary.

## Examples

### Plugin vsphere-iso
See the [examples folder](https://github.com/jetbrains-infra/packer-builder-vsphere/tree/master/examples) for examples of using `vsphere-iso`.

### Plugin vsphere-clone

#### Minimal Example
```json
{
  "builders": [
    {
      "type": "vsphere-clone",

      "vcenter_server": "vcenter.domain.com",
      "username": "root",
      "password": "secret",

      "template": "ubuntu",
      "vm_name":  "vm-1",
      "host":     "esxi-1.domain.com",

      "ssh_username": "root",
      "ssh_password": "secret"
    }
  ],
  "provisioners": [
    {
      "type": "shell",
      "inline": [ "echo hello" ]
    }
  ]
}
```

#### Complete Example
```json
{
  "variables": {
    "vsphere_password": "secret",
    "guest_password": "secret"
  },

  "builders": [
    {
      "type": "vsphere-clone",

      "vcenter_server": "vcenter.domain.com",
      "username": "root",
      "password": "{{user `vsphere_password`}}",
      "insecure_connection": true,
      "datacenter": "dc1",

      "template": "folder/ubuntu",
      "vm_name": "vm-1",
      "folder": "folder1/folder2",
      "host": "folder/esxi-1.domain.com",
      "resource_pool": "pool1/pool2",
      "datastore": "datastore1",
      "linked_clone": true,

      "CPUs": 2,
      "CPU_reservation": 1000,
      "CPU_limit": 2000,
      "RAM": 8192,
      "RAM_reservation": 2048,

      "ssh_username": "root",
      "ssh_password": "{{user `guest_password`}}",

      "shutdown_command": "echo '{{user `guest_password`}}' | sudo -S shutdown -P now",
      "shutdown_timeout": "5m",
      "create_snapshot": true,
      "convert_to_template": true
    }
  ],

  "provisioners": [
    {
      "type": "shell",
      "environment_vars": [
        "DEBIAN_FRONTEND=noninteractive"
      ],
      "execute_command": "echo '{{user `guest_password`}}' | {{.Vars}} sudo -ES bash -eux '{{.Path}}'",
      "inline": [
        "apt-get install -y zip"
      ]
    }
  ]
}
```

## Configuration Reference

### Specifying Clusters and Hosts
The `cluster` and `host` configuration options control where virtual machines will be created. This section applies to both the `vsphere-iso` and `vsphere-clone` builders.

#### ESXi Host Without Cluster
Only use the `host` option. Do not use the `cluster` option. Optionally specify a `resource_pool`.

```
"host": "esxi-2.vsphere65.test"
```

OR

```
"host": "esxi-2.vsphere65.test"
resource_pool": "pool1"
```

#### ESXi Cluster Without DRS
Use the `cluster` and `host `options.

```
"cluster": "cluster1",
"host": "esxi-2.vsphere65.test"
```

#### ESXi Cluster With DRS
Only use the `cluster` option. Do not use the `host` option. Optionally specify a `resource_pool`.

```
"cluster": "cluster2"
```

OR

```
"cluster": "cluster2"
resource_pool": "pool1"
```

### Plugin vsphere-iso

#### Required
* `vcenter_server`(string) - vCenter server hostname.
* `username`(string) - vSphere username.
* `password`(string) - vSphere password.
* `host`(string) - ESXi host where target VM is created. A full path must be specified if the host is in a folder. For example `folder/host`. See the `Specifying Clusters and Hosts` section above for more details.
* `cluster`(string)  - ESXi cluster where target VM is created. See the `Specifying Clusters and Host` section above for more details.
* `vm_name`(string) - Name of the new VM to create.
* `ssh_username`(string) - Username in guest OS.
* `ssh_password`(string) - Password to access guest OS. Only specify `ssh_password` or `ssh_private_key_file`, but not both.
* `ssh_private_key_file`(string) - Path to the SSH private key file to access guest OS. Only specify `ssh_password` or `ssh_private_key_file`, but not both.

#### Optional
* `boot_command`(array of strings) - List of commands to type when the VM is first booted. Used to initalize the operating system installer.
* `boot_order`(string) - Set VM boot order. Uses a comma delimiated string. Example ``"floppy,cdrom,ethernet,disk"``.
* `boot_wait`(string)  - Amount of time to wait for the VM to boot. Examples 45s and 10m. Defaults to 10s(10 seconds). See the Go Lang [ParseDuration](https://golang.org/pkg/time/#ParseDuration) documentation for full details.
* `convert_to_template`(boolean) - Convert VM to a template. Defaults to `false`.
* `CPUs`(number) - Number of CPU sockets.
* `CPU_limit`(number) - Upper limit of available CPU resources in MHz.
* `CPU_reservation`(number) - Amount of reserved CPU resources in MHz.
* `create_snapshot`(boolean) - Create a snapshot when set to `true`, so the VM can be used as a base for linked clones. Defaults to `false`.
* `datacenter`(string) - VMWare datacenter name. Required if there is more than one datacenter in vCenter.
* `datastore`(string) - VMWare datastore. Required if `host` is a cluster, or if `host` has multiple datastores.
* `disk_controller_type`(string) - Set VM disk controller type. Example `pvscsi`.
* `disk_size`(number) - The size of the disk in GB.
* `disk_thin_provisioned`(boolean) - Enable VMDK thin provisioning for VM. Defaults to `false`.
* `floppy_dirs`(array of strings) - Seems to not do anything useful yet. Not implemented.
* `floppy_files`(array of strings) - List of local files to be mounted to the VM floppy drive. Can be used to make Debian preseed or RHEL kickstart files available to the VM.
* `floppy_img_path`(string) - Data store path to a floppy image that will be mounted to the VM. Cannot be used with `floppy_files` or `floppy_dir` options. Example `[datastore1] ISO/VMware Tools/10.2.0/pvscsi-Windows8.flp`.
* `folder`(string) - VM folder to create the VM in.
* `guest_os_type`(string) - Set VM OS type. Defaults to `otherGuest`. See [here](https://pubs.vmware.com/vsphere-6-5/index.jsp?topic=%2Fcom.vmware.wssdk.apiref.doc%2Fvim.vm.GuestOsDescriptor.GuestOsIdentifier.html) for a full list of possible values.
* `insecure_connection`(boolean) - Do not validate vCenter server's TLS certificate. Defaults to `false`.
* `iso_paths`(array of strings) - List of data store paths to ISO files that will be mounted to the VM. Example `"[datastore1] ISO/ubuntu-16.04.3-server-amd64.iso"`.
* `network`(string) - Set network VM will be connected to.
* `network_card`(string) - Set VM network card type. Example `vmxnet3`.
* `RAM`(number) - Amount of RAM in MB.
* `RAM_reservation`(number) - Amount of reserved RAM in MB.
* `RAM_reserve_all`(boolean) - Reserve all available RAM. Defaults to `false`. Cannot be used together with `RAM_reservation`.
* `resource_pool`(string) - VMWare resource pool. Defaults to the root resource pool of the `host` or `cluster`.
* `shutdown_command`(string) - Specify a VM guest shutdown command. VMware guest tools are used by default.
* `shutdown_timeout`(string) - Amount of time to wait for graceful VM shutdown. Examples 45s and 10m. Defaults to 5m(5 minutes). See the Go Lang [ParseDuration](https://golang.org/pkg/time/#ParseDuration) documentation for full details.
* `usb_controller`(boolean) - Create US controller for virtual machine. Defaults to `false`.

### Plugin vsphere-clone

#### Required
* `vcenter_server`(string) - vCenter server hostname.
* `username`(string) - vSphere username.
* `password`(string) - vSphere password.
* `host`(string) - ESXi host where target VM is created. A full path must be specified if the host is in a folder. For example `folder/host`. See the `Specifying Clusters and Hosts` section above for more details.
* `cluster`(string)  - ESXi cluster where target VM is created. See the `Specifying Clusters and Host` section above for more details.
* `template`(string) - Name of source VM. Path is optional.
* `vm_name`(string) - Name of the new VM to create.
* `ssh_username`(string) - Username in guest OS.
* `ssh_password`(string) - Password to access guest OS. Only specify `ssh_password` or `ssh_private_key_file`, but not both.
* `ssh_private_key_file`(string) - Path to the SSH private key file to access guest OS. Only specify `ssh_password` or `ssh_private_key_file`, but not both.

#### Optional
* `boot_order`(string) Set the boot order of the VM. Uses a comma delimiated string. Example "floppy,cdrom,ethernet,disk".
* `boot_wait`(string)  Amount of time to wait for the VM to boot. Examples 45s and 10m. Defaults to 10s(10 seconds). See the Go Lang [ParseDuration](https://golang.org/pkg/time/#ParseDuration) documentation for full details.
* `convert_to_template`(boolean) - Convert VM to a template. Defaults to `false`.
* `CPUs`(number) - Number of CPU sockets. Inherited from `template` by default.
* `CPU_limit`(number) - Upper limit of available CPU resources in MHz. Inherited from `template` by default, set to `-1` for reset.
* `CPU_reservation`(number) - Amount of reserved CPU resources in MHz. Inherited from `template` by default.
* `create_snapshot`(boolean) - Create a snapshot when set to `true`, so the VM can be used as a base for linked clones. Defaults to `false`.
* `datacenter`(string) - VMWare datacenter name. Required if there is more than one datacenter in vCenter.
* `datastore`(string) - VMWare datastore. Required if `host` is a cluster, or if `host` has multiple datastores.
* `disk_size`(number) - The size of the disk in GB. Cannot be used together with `linked_clone`.
* `folder`(string) - VM folder to create the VM in.
* `insecure_connection`(boolean) - Do not validate vCenter server's TLS certificate. Defaults to `false`.
* `linked_clone`(boolean) - Create VM as a linked clone from latest snapshot. Defaults to `false`.
* `NestedHV`(boolean) - Enable nested hardware virtualization for VM. Defaults to `false`.
* `RAM`(number) - Amount of RAM in MB. Inherited from `template` by default.
* `RAM_reservation`(number) - Amount of reserved RAM in MB. Inherited from `template` by default.
* `RAM_reserve_all`(boolean) - Reserve all available RAM. Defaults to `false`. Cannot be used together with `RAM_reservation`.
* `resource_pool`(string) - VMWare resource pool. Defaults to the root resource pool of the `host` or `cluster`.
* `shutdown_command`(string) - Specify a VM guest shutdown command. VMware guest tools are used by default.
* `shutdown_timeout`(string) - Amount of time to wait for graceful VM shutdown. Examples 45s and 10m. Defaults to 5m(5 minutes). See the Go Lang [ParseDuration](https://golang.org/pkg/time/#ParseDuration) documentation for full details.
