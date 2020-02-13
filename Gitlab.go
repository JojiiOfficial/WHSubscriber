package main

import (
	"os"

	"github.com/JojiiOfficial/configor"
)

//GitlabActionStruct the struct for github webhooks
type GitlabActionStruct struct {
	Trigger string
	Filter  map[string]string
	EnvVars []string
	Actions []ActionItem
}

func createDefaultGitlabFile(file string) error {
	_, err := os.Stat(file)
	if err == nil {
		return nil
	}
	ghActionStruct := GitlabActionStruct{
		Trigger: "push",
		Filter:  map[string]string{"branch": "master"},
		EnvVars: []string{
			"PATH=/bin:/sbin:/usr/local/bin:/usr/local/sbin:/usr/sbin",
		},
		Actions: []ActionItem{
			ActionItem{
				Type:  ScriptActionItem,
				Value: "/some/script/Here",
			},
			ActionItem{
				Type:  CommandActionItem,
				Value: "cat /bin/bash",
			},
		},
	}
	_, err = configor.SetupConfig(&ghActionStruct, file, configor.NoChange)
	return err
}
