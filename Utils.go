package main

import (
	"log"
	"os"
	"path"
)

func inPortValid(port uint16) bool {
	return port > 0 && port < 65535
}

func getHome() string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln(err.Error())
		return ""
	}
	return home
}

func getDataPath() string {
	path := path.Join(getHome(), DataDir)
	s, err := os.Stat(path)
	if err != nil {
		err = os.Mkdir(path, 0770)
		if err != nil {
			log.Fatalln(err.Error())
		}
	} else if s != nil && !s.IsDir() {
		log.Fatalln("Datapath-name already taken by a file!")
	}
	return path
}
