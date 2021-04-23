package project

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/muidea/magicToolkit/core/util"
)

// Project core define
type Project struct {
	projectName  string
	rootPath     string
	templatePath string
}

// New new core
func New(rootPath string) *Project {
	return &Project{rootPath: rootPath}
}

func (s *Project) getFile(filePath string) (ret *os.File, err error) {
	existsFlag, existsErr := util.Exists(filePath)
	if !existsFlag || existsErr != nil {
		err = fmt.Errorf("file does not exist")
		return
	}

	ret, err = os.Open(filePath)
	return
}

func (s *Project) cloneTemplate(filePath string) (err error) {
	log.Print(filePath)
	templatePath := path.Join(s.templatePath, "project", filePath)
	fullPath := path.Join(s.rootPath, filePath)
	curPath, curName := path.Split(fullPath)
	curName = curName[:len(curName)-len(".template")]
	if curName == "main.go" {
		curPath = path.Join("cmd", s.projectName)
	}

	existsFlag, existsErr := util.Exists(curPath)
	if existsErr != nil {
		err = fmt.Errorf("unexpect, fileName:%s", curPath)
		return
	}
	if !existsFlag {
		err = os.MkdirAll(curPath, os.ModePerm)
		if err != nil {
			return
		}
	}

	filePath = path.Join(curPath, curName)
	existsFlag, existsErr = util.Exists(filePath)
	if existsErr != nil {
		err = fmt.Errorf("unexpect, fileName:%s", filePath)
		return
	}
	if existsFlag {
		return
	}

	fileHandle, fileErr := os.Create(filePath)
	if fileErr != nil {
		err = fileErr
		return
	}
	defer fileHandle.Close()

	templateHandle, templateErr := template.ParseFiles(templatePath)
	if templateErr != nil {
		err = templateErr
		return
	}

	data := map[string]string{"ProjectName": s.projectName}
	err = templateHandle.Execute(fileHandle, data)
	return
}

// Create create
func (s *Project) Create(projectName string) (err error) {
	s.projectName = projectName

	templatePath, templateErr := util.GetTemplatePath()
	if templateErr != nil {
		err = templateErr
		return
	}
	s.templatePath = templatePath

	err = s.walkTemplate("")
	return
}

func (s *Project) walkTemplate(projectPath string) (err error) {
	fullPath := path.Join(s.templatePath, "project", projectPath)
	existFlag, existErr := util.Exists(fullPath)
	if !existFlag || existErr != nil {
		err = fmt.Errorf("template file does not exist, projectPath:%s", projectPath)
		return
	}

	fdList, fdErr := ioutil.ReadDir(fullPath)
	if fdErr != nil {
		err = fdErr
		return
	}

	for _, fd := range fdList {
		if fd.IsDir() {
			err = s.walkTemplate(path.Join(projectPath, fd.Name()))
			if err != nil {
				return
			}

			continue
		}

		items := strings.Split(fd.Name(), ".")
		extName := items[len(items)-1]
		if extName != "template" {
			continue
		}

		err = s.cloneTemplate(path.Join(projectPath, fd.Name()))
		if err != nil {
			return
		}
	}

	return
}
