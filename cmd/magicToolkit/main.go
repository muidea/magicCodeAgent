package main

import (
	"fmt"
	"log"
	"os"

	"github.com/muidea/magicToolkit/core"
)

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
		}
	]
}
`

func main() {
	curPath, curErr := os.Getwd()
	if curErr != nil {
		log.Printf("get current work directory failed, err:%s", curErr.Error())
		return
	}

	core := &core.Core{}

	data := fmt.Sprintf(defaultTemplate, "TT", "TT")

	curErr = core.Create(curPath, []byte(data))
	if curErr != nil {
		log.Printf("Create template failed, err:%s", curErr.Error())
		return
	}
}
