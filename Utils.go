package main

import (
	"log"
	"os"
	"path"
	"strings"
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

func getCurrentDir() string {
	exec, err := os.Getwd()
	if err != nil {
		log.Fatalln(err.Error())
		return ""
	}
	return exec
}

func dirAbs(scriptPath string) (string, bool) {
	s, err := os.Stat(scriptPath)
	if err != nil || s == nil || !s.IsDir() {
		return scriptPath, false
	}
	if strings.HasPrefix(scriptPath, "/") {
		return scriptPath, true
	}

	if strings.HasPrefix(scriptPath, "./") {
		return path.Join(getCurrentDir(), scriptPath[2:]), true
	}

	if strings.HasPrefix(scriptPath, "~/") {
		return path.Join(getHome(), scriptPath[2:]), true
	}

	return path.Join(getCurrentDir(), scriptPath), true
}

func mapKeyByValue(val int8, m map[string]int8) string {
	for k, v := range m {
		if v == val {
			return k
		}
	}
	return ""
}

//FileExists reports whether the named file or directory exists.
func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
