package core

import (
	"encoding/json"
	"log"
	"os"
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

func (s *Core) walkTemplate(path string, template *model.Template) (err error) {
	fullPath := strings.Join([]string{path, template.Name}, "/")
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

func (s *Core) walkItem(path string, item *model.Item) (err error) {
	if item.IsDir() {
		fullPath := strings.Join([]string{path, item.Name}, "/")
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

	return
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return true, err
}
