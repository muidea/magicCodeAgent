package main

import (
	"fmt"
	"os"
)

func init() {
	actionList = append(actionList, &action{
		name:         "help",
		defaultValue: false,
		description:  "help information",
		funcPtr:      printUsage,
	})
}

var usage = `
  -action string
	createProject: create new project,and initialze
	updateDep: update project dependence
`

func printUsage(curPath string) {
	fmt.Printf("Usage of %s:", os.Args[0])
	fmt.Print(usage)
}
