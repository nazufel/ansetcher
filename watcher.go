package main

import (
	"bufio"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

var inventoryRootNotFound = errors.New("could not find provided inventory root")
var plainTextSecretsFound = errors.New("found plaintext secrets")

// watcher func boostraps the rest of the program
func watcher() error {

	var c conf

	c.getConfig("./ansible-secrets-watcher.yaml")

	err := checkForInventoryRoot(c.InventoryRoot)
	if err != nil {
		return err
	}

	secretFiles, err := directoryWalk(c.InventoryRoot)
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

type conf struct {
	InventoryRoot string
}

func (c *conf) getConfig(cf string) *conf {

	f, err := ioutil.ReadFile(cf)
	if err != nil {
		log.Fatal("ansible-secrets-watcher could not find config file")
	}

	err = yaml.Unmarshal(f, c)
	if err != nil {
		log.Fatalf("unmarshal error: %v", err)
	}

	return c
}

// checks to see if the inventory root passed exists, if not, it throws an error
func checkForInventoryRoot(dir string) error {

	// check that the inventory root file exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Println(inventoryRootNotFound)
		return inventoryRootNotFound
	}

	return nil
}

// directoryWalk walks directories from the root directory and looks for files that match a certain naming pattern. for now, that is "secrets.yml," but will be expanded later
func directoryWalk(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			if strings.Contains(info.Name(), "secrets.yml") {
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

// prints the error message when plaintext secrets are found
func printErrorMessage(plainTextSecretFiles []string) error {

	for pt := range plainTextSecretFiles {
		log.Printf("Error! Found Ansible Vault secrets file in plaintext during commit: %v. Please encrypt the file and reattempt to commit.", plainTextSecretFiles[pt])
	}
	return plainTextSecretsFound
}

//EOF
