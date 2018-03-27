package driver

import "testing"

const testDatastoreclusterName = "CHDC1EXP01"

func TestDatastoreclusterAcc(t *testing.T) {
	d := newTestDriver(t)
	dsc, err := d.FindDatastorecluster(testDatastoreclusterName)
	if err != nil {
		t.Fatalf("Cannot find specified datastorecluster: '%v'", testDatastoreclusterName)
	}
	info, err := dsc.Info("name")
	if err != nil {
		t.Fatalf("Cannot read datastorecluster properties: '%v'", err)
	}
	if info.Name != testDatastoreclusterName {
		t.Errorf("Wrong datastorecluster. expected: '%v', got: '%v' instead", testDatastoreclusterName, info.Name)
	}
}