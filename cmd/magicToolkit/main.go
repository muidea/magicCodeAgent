package main

import (
	"flag"
	"log"
	"os"
)

type action struct {
	name         string
	value        bool
	defaultValue bool
	description  string
	funcPtr      func(string)
}

var actionList = []*action{}
var help = false
var projectName = ""

func main() {
	for idx := range actionList {
		val := actionList[idx]
		flag.BoolVar(&val.value, val.name, val.defaultValue, val.description)
	}
	flag.StringVar(&projectName, "name", projectName, "project name")

	flag.Parse()
	curPath, curErr := os.Getwd()
	if curErr != nil {
		log.Printf("get current work directory failed, err:%s", curErr.Error())
		return
	}

	for _, val := range actionList {
		if val.value {
			val.funcPtr(curPath)
			return
		}
	}

	flag.Usage()
}
