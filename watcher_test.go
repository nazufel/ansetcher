package main

import (
	"log"
	"os"
	"testing"
)

func Test_checkForInventoryRoot(t *testing.T) {

	////////////
	// set up //
	////////////

	// make sure previous testing directories don't already exist
	if err := os.RemoveAll("./test-inventories"); err != nil {
		log.Fatal(err)
	}

	// create test inventory root
	if err := os.Mkdir("./test-inventories", 0777); err != nil {
		log.Fatal(err)
	}

	///////////////
	// test runs //
	///////////////

	testCases := []struct {
		name          string
		inventoryPath string
		want          error
	}{
		{"test existing inventory root", "./test-inventories", nil},
		// {"test non-existent inventory", "./bad-inventories/secrets.yml"},
	}

	// test the inventory root func
	for _, tc := range testCases {
		err := checkForInventoryRoot(tc.inventoryPath)
		if err != tc.want {
			t.Errorf("testcase %v, expected %v, got: %v", tc.name, tc.want, err)
		}
	}

	///////////////
	// tear down //
	///////////////

	// cleanup testing directories
	if err := os.RemoveAll("./test-inventories"); err != nil {
		log.Fatal(err)
	}
}
