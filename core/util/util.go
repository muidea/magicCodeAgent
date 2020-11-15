package util

import "os"

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
