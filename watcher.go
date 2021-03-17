// watcher.go

// Copyright 2020 Ryan Ross

// Functionality for Ansetcher

// ------------------------------------------------------------------------- //

package main

import (
	"bufio"
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var inventoryLocationVariableNotFound = errors.New("ansetcher: ANSIBLE_INVENTORIES_ROOT is a required environment variable. Please set an environment variable to define where the Ansible inventories directory is located relative to the root of this repository.")
var inventoryRootNotFound = errors.New("ansetcher: could not find provided inventory root")
var plainTextSecretsFound = errors.New("ansetcher: found plaintext secrets")
var secretsFileNameVariableNotFound = errors.New("ansetcher: ANSIBLE_SECRETS_FILE_NAME is a required environment variable. Please set an environment variable to define where the Ansible inventories directory is located in relative to the root of this repository.")

// conf holds configuration information such as where to find the Ansible inventories 
// directory and which secrets file naming convention to look for
type conf struct {
	// filesystem location of the Ansible inventories directory to walk 
	// and check for unencrypted Vault secrets
	InventoryRoot string

	// the secret files name to search for
	SecretFileName string
}

// watcher func boostraps the rest of the program
func watcher() error {

	var c conf

	err := c.getConfig()
	if err != nil {
		log.Fatal(err)
	}

	err = c.checkForInventoryRoot()
	if err != nil {
		return err
	}

	secretFiles, err := c.directoryWalk()
	if err != nil {
		log.Fatal(err)
	}

	plainTextAnsibleSecretFiles, err := findPlainTextAnsibleSecrets(secretFiles)

	if len(plainTextAnsibleSecretFiles) != 0 {
		err := printErrorMessage(plainTextAnsibleSecretFiles)
		return err
	}
	return nil
}

// checks to see if the inventory root passed exists, if not, it throws an error
func (c *conf) checkForInventoryRoot() error {

	// check that the inventory root file exists
	if _, err := os.Stat(c.InventoryRoot); os.IsNotExist(err) {
		log.Println(inventoryRootNotFound)
		log.Printf("ansetcher: searched in: %v but could not find inventory", c.InventoryRoot)
		return inventoryRootNotFound
	}

	return nil
}

// directoryWalk walks directories from the root directory and looks for files that match a certain naming pattern. for now, that is "secrets.yml," but will be expanded later
func (c *conf) directoryWalk() ([]string, error) {
	var files []string
	err := filepath.Walk(c.InventoryRoot, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			if strings.Contains(info.Name(), c.SecretFileName) {
				files = append(files, path)
			}
		}
		return nil
	})
	return files, err
}

func findPlainTextAnsibleSecrets(secretsFiles []string) ([]string, error) {

	var plainTextSecretFiles []string

	for sf := range secretsFiles {
		file, err := os.Open(secretsFiles[sf])
		if err != nil {
			log.Fatalf("failed to open %v: %s", file, err)
			return nil, err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)
		var textLines []string

		for scanner.Scan() {
			textLines = append(textLines, scanner.Text())
		}

		if !strings.Contains(textLines[0], "$ANSIBLE_VAULT;1.1;AES256") {
			plainTextSecretFiles = append(plainTextSecretFiles, secretsFiles[sf])
		}
	}

	return plainTextSecretFiles, nil
}

func (c *conf) getConfig() error {

	if len(os.Getenv("ANSIBLE_INVENTORIES_ROOT")) == 0 {
		return inventoryLocationVariableNotFound
	}

	c.InventoryRoot = os.Getenv("ANSIBLE_INVENTORIES_ROOT")

	if len(os.Getenv("ANSIBLE_SECRETS_FILE_NAME")) == 0 {
		return secretsFileNameVariableNotFound
	}

	c.SecretFileName = os.Getenv("ANSIBLE_SECRETS_FILE_NAME")
	return nil
}

// prints the error message when plaintext secrets are found
func printErrorMessage(plainTextSecretFiles []string) error {

	for pt := range plainTextSecretFiles {
		log.Printf("ansetcher: ERROR! Found Ansible Vault secrets file in plaintext during commit: %v. Please encrypt the file and reattempt to commit.", plainTextSecretFiles[pt])
	}
	return plainTextSecretsFound
}

//EOF
