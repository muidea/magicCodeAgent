package config

import "flag"

var dataFile = ""
var dataContent = ""

func init() {
	flag.StringVar(&dataFile, "dataFile", dataFile, "tool data file")
	flag.StringVar(&dataContent, "dataContent", dataContent, "tool data content")
}

func GetDataFile() string {
	return dataFile
}

func GetDataContent() string {
	return dataContent
}
