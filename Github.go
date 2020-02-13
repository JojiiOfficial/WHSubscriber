package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/JojiiOfficial/configor"
)

//GithubActionStruct the struct for github webhooks
type GithubActionStruct struct {
	Trigger string
	Filter  map[string]string
	EnvVars []string
	Actions []ActionItem
}

func createDefaultGithubFile(file string) error {
	_, err := os.Stat(file)
	if err == nil {
		return nil
	}
	ghActionStruct := GithubActionStruct{
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

//LoadGithubAction loads the action from a file
func LoadGithubAction(file string) (*GithubActionStruct, error) {
	action := GithubActionStruct{}
	err := configor.Load(&action, file)
	if err != nil {
		return nil, err
	}
	return &action, nil
}

//Run runs the github action
func (ghaction *GithubActionStruct) Run(payloadFile string) error {
	if len(ghaction.Actions) == 0 {
		return errors.New("no action defined")
	}
	fmt.Println(ghaction.Actions)
	return nil
}
