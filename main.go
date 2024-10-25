// Copyright 2024 The devsh authors

package main

import (
	"log"

	"github.com/kukushkin/devsh/cmd"
)

var (
	debug = "true" // indicates a debug build
)

func main() {
	log.SetFlags(0)
	if debug == "true" {
		log.SetFlags(log.Lshortfile)
	}

	cmd.Execute()
}
