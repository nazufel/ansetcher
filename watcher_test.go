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
		{"test non-existent inventory", "./bad-inventories/secrets.yml", inventoryRootNotFound},
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

func Test_walkDirectories(t *testing.T) {

	// TODO: clean up test dirs and files even if the test fail

	///////////
	// setup //
	///////////

	// clean up test-inventories dir for a clean slate to test on
	if err := os.RemoveAll("./test-inventories"); err != nil {
		log.Println(err)
	}

	// create test directories to walk
	testDirectories := [5]string{"./test-inventories",
		"./test-inventories/development",
		"./test-inventories/qa",
		"./test-inventories/stage",
		"./test-inventories/production"}

	for d := range testDirectories {
		if err := os.Mkdir(testDirectories[d], 0777); err != nil {
			log.Fatal(err)
		}
	}

	// create test files to walk
	testSecretsFiles := [4]string{"./test-inventories/development/secrets.yml",
		"./test-inventories/qa/secrets.yml",
		"./test-inventories/stage/secrets.yml",
		"./test-inventories/production/secrets.yml"}

	// create test files to walk
	testDefaultsFiles := [4]string{"./test-inventories/development/defaults.yml",
		"./test-inventories/qa/defaults.yml",
		"./test-inventories/stage/defaults.yml",
		"./test-inventories/production/defaults.yml"}

	for f := range testSecretsFiles {
		if _, err := os.Create(testSecretsFiles[f]); err != nil {
			log.Fatal(err)
		}
	}

	for f := range testDefaultsFiles {
		if _, err := os.Create(testDefaultsFiles[f]); err != nil {
			log.Fatal(err)
		}
	}

	///////////////
	// run tests //
	///////////////

	inventoryRoot := "./test-inventories"

	walkedFiles, err := directoryWalk(inventoryRoot)
	if err != nil {
		t.Error(err)
	}

	// test if the legnth of walked files is expected
	if len(walkedFiles) != 4 {
		// clean up test-inventories dir for a clean slate to test on
		if err := os.RemoveAll("./test-inventories"); err != nil {
			log.Fatal(err)
		}
		t.Errorf("expected to find: %v secrets files, found: %v", 4, len(walkedFiles))
	}

	// test if all the found files are expected
	for wf := range walkedFiles {
		if walkedFiles[wf] != testSecretsFiles[wf] {
			t.Errorf("expected to find secret file %v, found %v", testSecretsFiles[wf], walkedFiles[wf])
		}
	}

	//////////////
	// teardown //
	//////////////

	// clean up test-inventories dir for a clean slate to test on
	if err := os.RemoveAll("./test-inventories"); err != nil {
		log.Fatal(err)
	}
}

//EOF
