package driver

import "testing"

func TestFolderAcc(t *testing.T) {
	d := newTestDriver(t)
	f, err := d.FindFolder("/esxi-1.vsphere65.test/folder1/folder2")
	if err != nil {
		t.Fatalf("Cannot find the default folder '%v': %v", "/esxi-1.vsphere65.test/folder1/folder2", err)
	}
	path, err := f.Path()
	if err != nil {
		t.Fatalf("Cannot read folder name: %v", err)
	}
	if path != "/esxi-1.vsphere65.test/folder1/folder2" {
		t.Errorf("Wrong folder. expected: '/esxi-1.vsphere65.test/folder1/folder2', got: '%v'", path)
	}
}
