package project

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"text/template"
)

// Core core define
type Core struct {
	projectName  string
	rootPath     string
	templatePath string
}

// NewCore new core
func NewCore(rootPath string) *Core {
	return &Core{rootPath: rootPath}
}

func (s *Core) getFile(filePath string) (ret *os.File, err error) {
	existsFlag, existsErr := exists(filePath)
	if !existsFlag || existsErr != nil {
		err = fmt.Errorf("file does not exist")
		return
	}

	ret, err = os.Open(filePath)
	return
}

func (s *Core) getTemplatePath() (ret string, err error) {
	toolKitPath := os.Getenv("MAGICTOOLKIT_PATH")
	if toolKitPath == "" {
		err = fmt.Errorf("no found toolkitPath")
		return
	}

	templatePath := path.Join(toolKitPath, "template")
	existsFlag, existsErr := exists(templatePath)
	if !existsFlag || existsErr != nil {
		err = fmt.Errorf("toolkit template does not exist")
		return
	}

	ret = templatePath
	return
}

func (s *Core) cloneTemplate(filePath string) (err error) {
	templatePath := path.Join(s.templatePath, filePath)
	fullPath := path.Join(s.rootPath, filePath)
	curPath, curName := path.Split(fullPath)
	existsFlag, existsErr := exists(curPath)
	if existsFlag || existsErr != nil {
		err = fmt.Errorf("unexpect, fileName:%s", filePath)
		return
	}

	err = os.MkdirAll(curPath, os.ModePerm)
	if err != nil {
		return
	}

	curName = curName[:len(curName)-len(".template")]
	fileHandle, fileErr := os.Create(path.Join(curPath, curName))
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
func (s *Core) Create(projectName string) (err error) {
	curPath := s.rootPath
	s.projectName = projectName

	templatePath, templateErr := s.getTemplatePath()
	if templateErr != nil {
		err = templateErr
		return
	}
	s.templatePath = templatePath

	err = s.walkTemplate(curPath)
	return
}

func (s *Core) walkTemplate(projectPath string) (err error) {
	fullPath := path.Join(s.templatePath, projectPath)
	existFlag, existErr := exists(fullPath)
	if !existFlag || existErr != nil {
		err = fmt.Errorf("template file does not exist")
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

		err = s.cloneTemplate(path.Join(projectPath, fd.Name()))
		if err != nil {
			return
		}
	}

	return
}

func exists(itemPath string) (bool, error) {
	_, err := os.Stat(itemPath)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return true, err
}
