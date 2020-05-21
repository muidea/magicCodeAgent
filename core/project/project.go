package project

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"

	"github.com/muidea/magicToolkit/model"
)

// Core core define
type Core struct {
	rootPath string
}

// NewCore new core
func NewCore(rootPath string) *Core {
	return &Core{rootPath: rootPath}
}

func (s *Core) getTemplateFile(templateFile string) (ret *os.File, err error) {
	toolKitPath := os.Getenv("MAGICTOOLKIT_PATH")
	if toolKitPath == "" {
		err = fmt.Errorf("no found toolkitPath")
		return
	}

	templatePath := path.Join(toolKitPath, "template")
	existsFlag, existsErr := exists(templatePath)
	if !existsFlag || existsErr != nil {
		err = fmt.Errorf("no found toolkit template path")
		return
	}

	templateFile = path.Join(templatePath, fmt.Sprintf("%s.template", templateFile))
	_, fileName := path.Split(templateFile)
	if fileName == "main.go.template" {
		templateFile = path.Join(templatePath, fileName)
	}

	existsFlag, existsErr = exists(templateFile)
	if !existsFlag || existsErr != nil {
		err = fmt.Errorf("no found toolkit template file")
		return
	}

	ret, err = os.Open(templateFile)

	return
}

func (s *Core) cloneTemplateFile(fullFilePath string) (err error) {
	projectName := os.Getenv("projectName")
	if projectName == "" {
		err = fmt.Errorf("illegal projectName")
		return
	}

	filePtr, fileErr := os.Create(fullFilePath)
	if fileErr != nil {
		err = fileErr
		return
	}
	defer filePtr.Close()

	fileName := fullFilePath[len(s.rootPath)+len(projectName)+2:]
	templateFile, templateErr := s.getTemplateFile(fileName)
	if templateErr == nil {
		defer templateFile.Close()

		buf := bufio.NewReader(templateFile)
		for {
			strLine, strErr := buf.ReadString('\n')
			if strErr != nil {
				if strErr == io.EOF {
					err = nil
					return
				}

				err = strErr
				return
			}

			strLine = strings.ReplaceAll(strLine, "%projectName%", projectName)

			_, err = filePtr.WriteString(strLine)
			if err != nil {
				return
			}
		}
	}

	return
}

// Create create
func (s *Core) Create(data []byte) (err error) {
	template := &model.Template{}

	err = json.Unmarshal(data, template)
	if err != nil {
		return
	}

	curPath := s.rootPath
	err = s.walkTemplate(curPath, template)
	return
}

func (s *Core) walkTemplate(itemPath string, template *model.Template) (err error) {
	fullPath := strings.Join([]string{itemPath, template.Name}, "/")
	existFlag, existErr := exists(fullPath)
	if existErr != nil {
		err = existErr
		return
	}

	if !existFlag {
		err = os.MkdirAll(fullPath, os.ModePerm)
		if err != nil {
			log.Printf("mkdir failed, err:%s", err.Error())
			return
		}
	}

	for _, val := range template.Items {
		err = s.walkItem(fullPath, val)
		if err != nil {
			return
		}
	}

	return
}

func (s *Core) walkItem(itemPath string, item *model.Item) (err error) {
	itemPath = strings.ToLower(itemPath)
	if item.IsDir() {
		fullPath := strings.Join([]string{itemPath, item.Name}, "/")
		existFlag, existErr := exists(fullPath)
		if existErr != nil {
			err = existErr
			return
		}

		if !existFlag {
			err = os.MkdirAll(fullPath, os.ModePerm)
			if err != nil {
				log.Printf("mkdir failed, err:%s", err.Error())
				return
			}
		}

		for _, val := range item.SubItems {
			err = s.walkItem(fullPath, val)
			if err != nil {
				return
			}
		}

		return
	}

	err = s.createFile(itemPath, item)
	return
}

func (s *Core) createFile(itemPath string, item *model.Item) (err error) {
	fullPath := strings.Join([]string{itemPath, item.Name}, "/")
	existFlag, existErr := exists(fullPath)
	if existErr != nil {
		err = existErr
		return
	}

	if existFlag {
		return
	}

	err = s.cloneTemplateFile(fullPath)
	if err == nil {
		return
	}

	filePtr, fileErr := os.Create(fullPath)
	if fileErr != nil {
		err = fileErr
		return
	}
	defer filePtr.Close()

	fileExt := ""
	itemsInfo := strings.Split(item.Name, ".")
	if len(itemsInfo) > 1 {
		fileExt = itemsInfo[len(itemsInfo)-1]
	}
	switch fileExt {
	case "go":
		_, pathName := path.Split(itemPath)
		if item.Name == "main.go" {
			pathName = "main"
		}

		filePtr.WriteString(fmt.Sprintf("package %s", pathName))
	}
	return nil
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
