package main

import (
	"log"
	"os"
	"path/filepath"
)

InventoryRootNotFound := errors.New("could not find provided inventory root")

// watcher func boostraps the rest of the program
func watcher() error {

	var err error
	// root := "inventories"

	// rootExists, err := os.Stat(root)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// if os.IsNotExist(err) {

	// }

	// walkedFiles, err := walkDirectories(root)
	// if err != nil {
	// 	log.Println(err)
	// }

	// for i := 0; i < len(walkedFiles); i++ {
	// 	fmt.Println(walkedFiles[i])
	// }
	return err
}

// checks to see if the inventory root passed exists, if not, it throws an error
func checkForInventoryRoot(dir string) error {

	// check that the inventory root file exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Println(err)
		return err
	}

	return nil
}

func walkDirectories(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}
