package main

import "os"

func main() {

	err := watcher()
	if err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
