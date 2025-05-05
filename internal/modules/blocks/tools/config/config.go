package config

import "flag"

var dataFile = ""

func init() {
	flag.StringVar(&dataFile, "dataFile", dataFile, "tool data file")
}

func GetDataFile() string {
	return dataFile
}
