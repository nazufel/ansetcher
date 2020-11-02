package main

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func Test_checkForInventoryRoot(t *testing.T) {

	////////////
	// set up //
	////////////

	// make sure previous testing directories don't already exist
	if err := os.RemoveAll("test-inventories"); err != nil {
		log.Fatal(err)
	}

	// create test inventory root
	if err := os.Mkdir("test-inventories", 0777); err != nil {
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
		{"test existing inventory root", "test-inventories", nil},
		{"test non-existent inventory", "bad-inventories/secrets.yml", inventoryRootNotFound},
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
	if err := os.RemoveAll("test-inventories"); err != nil {
		log.Fatal(err)
	}
}

func Test_walkDirectories(t *testing.T) {

	// TODO: clean up test dirs and files even if the test fail

	///////////
	// setup //
	///////////

	// clean up test-inventories dir for a clean slate to test on
	if err := os.RemoveAll("test-inventories"); err != nil {
		log.Println(err)
	}

	// create test directories to walk
	testDirectories := [5]string{"test-inventories",
		"test-inventories/development",
		"test-inventories/qa",
		"test-inventories/stage",
		"test-inventories/production"}

	for d := range testDirectories {
		if err := os.Mkdir(testDirectories[d], 0777); err != nil {
			log.Fatal(err)
		}
	}

	// create test files to walk
	testSecretsFiles := [4]string{"test-inventories/development/secrets.yml",
		"test-inventories/production/secrets.yml",
		"test-inventories/qa/secrets.yml",
		"test-inventories/stage/secrets.yml"}

	// create test files to walk
	testDefaultsFiles := [4]string{"test-inventories/development/defaults.yml",
		"test-inventories/production/defaults.yml",
		"test-inventories/qa/defaults.yml",
		"test-inventories/stage/defaults.yml"}

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

	inventoryRoot := "test-inventories"

	walkedFiles, err := directoryWalk(inventoryRoot)
	if err != nil {
		t.Error(err)
	}

	// test if the legnth of walked files is expected
	if len(walkedFiles) != 4 {
		// clean up test-inventories dir for a clean slate to test on
		if err := os.RemoveAll("test-inventories"); err != nil {
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
	if err := os.RemoveAll("test-inventories"); err != nil {
		log.Fatal(err)
	}
}

func Test_findAnsibleSecrets(t *testing.T) {

	///////////
	// setup //
	///////////

	// clean up test-inventories dir for a clean slate to test on
	if err := os.RemoveAll("test-inventories"); err != nil {
		log.Println(err)
	}

	// create test directories to walk
	testDirectories := [5]string{"test-inventories",
		"test-inventories/development",
		"test-inventories/qa",
		"test-inventories/stage",
		"test-inventories/production"}

	for d := range testDirectories {
		if err := os.Mkdir(testDirectories[d], 0777); err != nil {
			log.Fatal(err)
		}
	}

	// create test files to walk
	testSecretFiles := []struct {
		filePath    string
		fileContent []byte
		permissions os.FileMode
	}{
		{"test-inventories/development/secrets.yml", []byte("$ANSIBLE_VAULT;1.1;AE256"), 0777},
		{"test-inventories/production/secrets.yml", []byte("$ANSIBLE_VAULT;1.1;AE256"), 0777},
		{"test-inventories/qa/secrets.yml", []byte("$ANSIBLE_VAULT;1.1;AE256"), 0777},
		{"test-inventories/stage/secrets.yml", []byte("$ANSIBLE_VAULT;1.1;AE256"), 0777},
	}

	// create test files to walk
	testDefaultFiles := []struct {
		filePath    string
		fileContent []byte
		permissions os.FileMode
	}{
		{"test-inventories/development/defaults.yml", []byte("not a secret"), 0777},
		{"test-inventories/production/defaults.yml", []byte("not a secret"), 0777},
		{"test-inventories/qa/defaults.yml", []byte("not a secret"), 0777},
		{"test-inventories/stage/defaults.yml", []byte("not a secret"), 0777},
	}

	var testSecretFilePaths []string
	var testDefaultFilePaths []string

	for _, tf := range testSecretFiles {
		if err := ioutil.WriteFile(tf.filePath, tf.fileContent, tf.permissions); err != nil {
			log.Fatal(err)
		}
		testSecretFilePaths = append(testSecretFilePaths, tf.filePath)
	}
	for _, tf := range testDefaultFiles {
		if err := ioutil.WriteFile(tf.filePath, tf.fileContent, tf.permissions); err != nil {
			log.Fatal(err)
		}
		testDefaultFilePaths = append(testDefaultFilePaths, tf.filePath)
	}

	///////////////
	// run tests //
	///////////////

	plainTextSecretFiles, err := findPlainTextAnsibleSecrets(testSecretFilePaths)
	if err != nil {
		log.Fatal(err)
	}

	// test 0: returned text paths should equal a certain number
	if len(plainTextSecretFiles) != 0 {
		t.Errorf("expected to find 0 plain text ansible secrets files, found %v", len(plainTextSecretFiles))
	}

	// test 1: returned text paths should have certain files
	for rf := range plainTextSecretFiles {
		if plainTextSecretFiles[rf] != testSecretFilePaths[rf] {
			t.Errorf("expected returned file: %v to equal secret test file: %v", plainTextSecretFiles[rf], testSecretFilePaths[rf])
		}
	}

	// test 2: returned text paths should not have certain files
	for rf := range plainTextSecretFiles {
		if plainTextSecretFiles[rf] == testDefaultFilePaths[rf] {
			t.Errorf("expected returned file: %v to not equal default test file: %v", plainTextSecretFiles[rf], testDefaultFilePaths[rf])
		}
	}
	//////////////
	// teardown //
	//////////////

	// clean up test-inventories dir for a clean slate to test on
	if err := os.RemoveAll("test-inventories"); err != nil {
		log.Fatal(err)
	}
}

//EOF
