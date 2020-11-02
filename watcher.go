package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var inventoryRootNotFound = errors.New("could not find provided inventory root")

// watcher func boostraps the rest of the program
func watcher() error {

	inventoryRoot := "./inventories"

	err := checkForInventoryRoot(inventoryRoot)
	if err != nil {
		return err
	}

	secretFiles, err := directoryWalk(inventoryRoot)
	if err != nil {
		log.Fatal(err)
	}

	for sf := range secretFiles {
		fmt.Println(secretFiles[sf])
	}

	return err
}

// checks to see if the inventory root passed exists, if not, it throws an error
func checkForInventoryRoot(dir string) error {

	// check that the inventory root file exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Println(inventoryRootNotFound)
		return inventoryRootNotFound
	}

	log.Printf("found %v, proceeding", dir)
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
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)
		var textLines []string

		for scanner.Scan() {
			textLines = append(textLines, scanner.Text())
		}

		if !strings.Contains(textLines[0], "$ANSIBLE_VAULT;1.1;AE265") {
			plainTextSecretFiles = append(plainTextSecretFiles, secretsFiles[sf])
		}
	}

	log.Println(plainTextSecretFiles)

	return plainTextSecretFiles, nil
}
