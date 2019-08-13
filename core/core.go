package core

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/muidea/magicToolkit/model"
)

// Core core define
type Core struct {
}

// Create create
func (s *Core) Create(curPath string, data []byte) (err error) {
	template := &model.Template{}

	err = json.Unmarshal(data, template)
	if err != nil {
		return
	}

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

	err = createFile(itemPath, item)
	return
}

func createFile(itemPath string, item *model.Item) (err error) {
	fullPath := strings.Join([]string{itemPath, item.Name}, "/")
	existFlag, existErr := exists(fullPath)
	if existErr != nil {
		err = existErr
		return
	}

	if existFlag {
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
