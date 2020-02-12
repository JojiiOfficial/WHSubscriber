package main

import (
	"os"

	"github.com/JojiiOfficial/configor"
)

//GitlabActionStruct the struct for github webhooks
type GitlabActionStruct struct {
	Trigger string
	Filter  map[string]string
	Actions []string
}

func createDefaultGitlabFile(file string) error {
	_, err := os.Stat(file)
	if err == nil {
		return nil
	}
	ghActionStruct := GithubActionStruct{
		Trigger: "push",
		Filter:  map[string]string{"branch": "master"},
		Actions: []string{"/a/script/here/to/run.sh"},
	}
	_, err = configor.SetupConfig(&ghActionStruct, file, configor.NoChange)
	return err
}
