package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/muidea/magicToolkit/core"
)

var projectName = ""

func init() {
	actionList = append(actionList, &action{
		name:         "createProject",
		defaultValue: false,
		description:  "create new project,and initialze",
		funcPtr:      createProject,
	})

	flag.StringVar(&projectName, "name", projectName, "project name")
}

var defaultTemplate = `
{
	"name": "%s",
	"items": [
		{
			"name": "cmd",
			"type": 0,
			"subItems": [
				{
					"name": "%s",
					"type": 0,
					"subItems": [
						{
							"name": "main.go",
							"type": 1,
							"subItems": []
						}
					]
				}
			]
		},
		{
			"name": "core",
			"type": 0,
			"subItems": [
				{
					"name": "core.go",
					"type": 1,
					"subItems": []
				},
				{
					"name": "route",
					"type": 0,
					"subItems": []
				},
				{
					"name": "handler",
					"type": 0,
					"subItems": []
				}								
			]
		},
		{
			"name": "model",
			"type": 0,
			"subItems": []
		},
		{
			"name": "config",
			"type": 0,
			"subItems": [
				{
					"name": "config.go",
					"type": 1,
					"subItems": []
				}
			]
		},
		{
			"name": "docker",
			"type": 0,
			"subItems": [
				{
					"name": "bootstrap.sh",
					"type": 1,
					"subItems": []
				},{
					"name": "build.sh",
					"type": 1,
					"subItems": []
				},{
					"name": "Dockerfile",
					"type": 1,
					"subItems": []
				},{
					"name": "wait-for-it.sh",
					"type": 1,
					"subItems": []
				}
			]
		}
	]
}
`

func createProject(curPath string) {
	if projectName == "" {
		fmt.Printf("must be input project name")
		return
	}

	core := &core.Core{}

	data := fmt.Sprintf(defaultTemplate, projectName, projectName)

	curErr := core.Create(curPath, []byte(data))
	if curErr != nil {
		log.Printf("Create template failed, err:%s", curErr.Error())
		return
	}
}
