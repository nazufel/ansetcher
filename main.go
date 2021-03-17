// main.go

// Copyright 2020 Ryan Ross

// Main function for the Ansetcher

// ------------------------------------------------------------------------- //

package main

import "os"

func main() {

	err := watcher()
	if err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
