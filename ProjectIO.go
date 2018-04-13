package main

import (
	"os"
)

var dbLocation string

func init() {
	dbLocation = os.Getenv("XDG_CONFIG_HOME")

	if dbLocation == "" {
		dbLocation = os.Getenv("HOME") + "/.config"
	}

	dbLocation += "/goplan"

	os.MkdirAll(dbLocation, os.ModePerm)
}

func SaveFile(p []*Project) {

	file, _ := os.Create(dbLocation + "/database.json")
	defer file.Close()

	SaveBatch(p, file)
}

func LoadFile() []*Project {

	file, _ := os.Open(dbLocation + "/database.json")
	defer file.Close()

	projs, _ := LoadBatch(file)

	return projs
}
