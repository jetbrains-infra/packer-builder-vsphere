package main

const BuilderId = "LizaTretyakova.vsphere"

type ArtifactFile struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
}

type Artifact struct {
	VMName        string         `json:"vm_name"`
}

func (a *Artifact) BuilderId() string {
	return BuilderId
}

func (a *Artifact) Files() []string {
	return []string{}
}

func (a *Artifact) Id() string {
	return a.VMName
}

func (a *Artifact) String() string {
	return a.VMName
}

func (a *Artifact) State(name string) interface{} {
	return nil
}

func (a *Artifact) Destroy() error {
	return nil
}
