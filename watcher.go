package main

import (
	"errors"
	"log"
	"os"
	"path/filepath"
)

var inventoryRootNotFound = errors.New("could not find provided inventory root")

// watcher func boostraps the rest of the program
func watcher() error {

	inventoryRoot := "./inventories"

	err := checkForInventoryRoot(inventoryRoot)
	if err != nil {
		return err
	}

	return err
}

// checks to see if the inventory root passed exists, if not, it throws an error
func checkForInventoryRoot(dir string) error {

	// check that the inventory root file exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Println(err)
		return inventoryRootNotFound
	}

	return nil
}

// directoryWalk walks directories from the root directory and looks for files that match a certain naming pattern. for now, that is "secrets.yml," but will be expanded later
func directoryWalk(root string) ([]string, error) {
	// TODO: stopped here. have a solid unit test. now need to get this func working
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}
