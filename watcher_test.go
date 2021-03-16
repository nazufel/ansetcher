// watcher_test.go

// Copyright 2020 Ryan Ross

// Test cases for the Ansible Secrets Watcher

// ------------------------------------------------------------------------- //

package main

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func (c *conf) Test_checkForInventoryRoot(t *testing.T) {

	////////////
	// set up //
	////////////


	c.InventoryRoot = "test-inventories"

	// make sure previous testing directories don't already exist
	if err := os.RemoveAll(c.InventoryRoot); err != nil {
		log.Fatal(err)
	}

	// create test inventory root
	if err := os.Mkdir(c.InventoryRoot, 0644); err != nil {
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
		{"test existing inventory root", c.InventoryRoot, nil},
		{"test non-existent inventory", "bad-inventories", inventoryRootNotFound},
	}

	// test the inventory root func
	for _, tc := range testCases {
		err := c.checkForInventoryRoot()
		if err != tc.want {
			t.Errorf("testcase %v, expected %v, got: %v", tc.name, tc.want, err)
		}
	}

	///////////////
	// tear down //
	///////////////

	// cleanup testing directories
	if err := os.RemoveAll(c.InventoryRoot); err != nil {
		log.Fatal(err)
	}
}

func (c *conf) Test_walkDirectories(t *testing.T) {

	// TODO: clean up test dirs and files even if the test fail

	///////////
	// setup //
	///////////


	c.InventoryRoot = "test-inventories"
	c.SecretFileName = "secrets.yaml"

	// clean up test-inventories dir for a clean slate to test on
	if err := os.RemoveAll(c.InventoryRoot); err != nil {
		log.Println(err)
	}

	// create test directories to walk
	testDirectories := [5]string{c.InventoryRoot,
		c.InventoryRoot + "/development",
		c.InventoryRoot + "/qa",
		c.InventoryRoot + "/stage",
		c.InventoryRoot + "/production"}

	for d := range testDirectories {
		if err := os.Mkdir(testDirectories[d], 0644); err != nil {
			log.Fatal(err)
		}
	}

	// create test files to walk
	testSecretsFiles := [4]string{c.InventoryRoot + "/development/" + c.SecretFileName,
		c.InventoryRoot + "/production/" + c.SecretFileName,
		c.InventoryRoot + "/qa/" + c.SecretFileName,
		c.InventoryRoot + "/stage/" + c.SecretFileName}

	// create test files to walk
	testDefaultsFiles := [4]string{c.InventoryRoot + "/development/defaults.yml",
		c.InventoryRoot + "/production/defaults.yml",
		c.InventoryRoot + "/qa/defaults.yml",
		c.InventoryRoot + "/stage/defaults.yml"}

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

	walkedFiles, err := c.directoryWalk()
	if err != nil {
		t.Error(err)
	}

	// test if the legnth of walked files is expected
	if len(walkedFiles) != 4 {
		// clean up test-inventories dir for a clean slate to test on
		if err := os.RemoveAll(c.InventoryRoot); err != nil {
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
	if err := os.RemoveAll(c.InventoryRoot); err != nil {
		log.Fatal(err)
	}
}

func (c *conf) Test_findAnsibleSecrets(t *testing.T) {

	///////////
	// setup //
	///////////

	// clean up test-inventories dir for a clean slate to test on
	if err := os.RemoveAll(c.InventoryRoot); err != nil {
		log.Println(err)
	}

	// create test directories to walk
	testDirectories := [5]string{c.InventoryRoot,
		c.InventoryRoot + "/development",
		c.InventoryRoot + "/qa",
		c.InventoryRoot + "/stage",
		c.InventoryRoot + "/production"}

	for d := range testDirectories {
		if err := os.Mkdir(testDirectories[d], 0644); err != nil {
			log.Fatal(err)
		}
	}

	// create test files to walk
	testSecretFiles := []struct {
		filePath    string
		fileContent []byte
		permissions os.FileMode
	}{
		{c.InventoryRoot + "/development/" + c.SecretFileName, []byte("plaintext secret. not good!"), 0644},
		{c.InventoryRoot + "/production/" + c.SecretFileName, []byte("plaintext secret. not good!"), 0644},
		{c.InventoryRoot + "/qa/" + c.SecretFileName, []byte("$ANSIBLE_VAULT;1.1;AES256\n12341234123412341\n1234"), 0644},
		{c.InventoryRoot + "/stage/" + c.SecretFileName, []byte("$ANSIBLE_VAULT;1.1;AES256\n12341234123432\n1234"), 0644},
	}

	// create test files to walk
	testDefaultFiles := []struct {
		filePath    string
		fileContent []byte
		permissions os.FileMode
	}{
		{c.InventoryRoot + "/development/" + c.SecretFileName, []byte("not a secret"), 0644},
		{c.InventoryRoot + "/production/" + c.SecretFileName, []byte("not a secret"), 0644},
		{c.InventoryRoot + "/qa/" + c.SecretFileName, []byte("not a secret"), 0644},
		{c.InventoryRoot + "/stage/" + c.SecretFileName, []byte("not a secret"), 0644},
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

	plainTextSecretFiles, _ := findPlainTextAnsibleSecrets(testSecretFilePaths)

	// test 0: returned text paths should equal a certain number
	if len(plainTextSecretFiles) != 2 {
		t.Errorf("expected to find 2 plain text ansible secrets files, found %v", len(plainTextSecretFiles))
	}

	// test 1: returned text paths should have certain files
	for rf := range plainTextSecretFiles {
		// skip over the first and third plainTextSecretFiles since they are actually plain text
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
	if err := os.RemoveAll(c.InventoryRoot); err != nil {
		log.Fatal(err)
	}
}

func Test_getConfig(t *testing.T) {
	
	var c conf

	// test happy path with set variables

	// set up the variables
	os.Setenv("ANSIBLE_INVENTORIES_ROOT", "./inventories")
	os.Setenv("ANSIBLE_SECRETS_FILE_NAME", "secrets.yaml")

	err := c.getConfig()
	if err != nil {
		t.Errorf("expected to find environment variables, but did not")
	}

	// test unhappy path with variables unset

	// test for not finding the inventory root
	os.Unsetenv("ANSIBLE_INVENTORIES_ROOT")
	err = c.getConfig()
	if err != inventoryLocationVariableNotFound {
		t.Errorf("expected to be unable to find environment variable: ANSIBLE_INVENTORIES_ROOT, but it was found.")
	}

	// test for not finding the secrets file name
	os.Setenv("ANSIBLE_INVENTORIES_ROOT", "./inventories")
	os.Unsetenv("ANSIBLE_SECRETS_FILE_NAME")
	err = c.getConfig()
	if err != secretsFileNameVariableNotFound {
		t.Errorf("expected to be unable to find environment variable: ANSIBLE_SECRETS_FILE_NAME, but it was found.")
	}

}
//EOF
