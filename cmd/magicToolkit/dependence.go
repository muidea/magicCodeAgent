package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func init() {
	actionList = append(actionList, &action{
		name:         "updateDep",
		defaultValue: false,
		description:  "update project dependence",
		funcPtr:      updateDep,
	})
}

func filterPath(pathVal string) (ret bool) {
	inGitControl := false
	inDepFlag := false
	inVendor := false

	files, _ := ioutil.ReadDir(pathVal)
	for _, fi := range files {
		fileName := fi.Name()
		if fi.IsDir() {
			if fileName == ".git" {
				inGitControl = true
				continue
			}

			if fileName == "vendor" {
				inVendor = true
				continue
			}

			continue
		} else {
			if fileName == "Gopkg.toml" {
				inDepFlag = true
				continue
			}
		}
	}

	return inGitControl && inDepFlag && inVendor
}

func updateDep(curPath string) {
	prePath, preErr := os.Getwd()
	if preErr != nil {
		log.Printf("get current work path failed, err:%s", preErr.Error())
		return
	}

	filepath.Walk(curPath, func(subPath string, subInfo os.FileInfo, subErr error) error {
		if !subInfo.IsDir() {
			return nil
		}
		if subInfo.Name() == ".git" {
			return nil
		}

		if filterPath(subPath) {
			err := os.Chdir(subPath)
			if err == nil {
				cmd := exec.Command("dep", "ensure", "-update")
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				err = cmd.Run()
			}
			if err == nil {
				log.Printf("----%s success", subPath)
			} else {
				log.Printf("----%s failed", subPath)
			}
		}

		return nil
	})

	os.Chdir(prePath)
}
