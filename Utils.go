package main

import (
	"log"
	"os"
	"path"

	gaw "github.com/JojiiOfficial/GoAw"
)

func getDataPath() string {
	path := path.Join(gaw.GetHome(), DataDir)
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

func mapKeyByValue(val int8, m map[string]int8) string {
	for k, v := range m {
		if v == val {
			return k
		}
	}
	return ""
}
