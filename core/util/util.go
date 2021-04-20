package util

import (
	"os"
	"strings"

	"github.com/muidea/magicOrm/provider/remote"
)

func Exists(itemPath string) (bool, error) {
	_, err := os.Stat(itemPath)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return true, err
}

func GetModelName(modelPtr *remote.Object) string {
	items := strings.Split(modelPtr.GetName(), ".")
	return items[1]
}
