package main

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/JojiiOfficial/configor"
)

//GithubActionStruct the struct for github webhooks
type GithubActionStruct struct {
	Trigger string
	Filter  map[string]string
	EnVars  []string
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
		EnVars: []string{
			"PATH=/bin/",
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
	for _, action := range ghaction.Actions {
		switch action.Type {
		case ScriptActionItem:
			{

			}
		case CommandActionItem:
			{
				if len(strings.Trim(action.Value, " ")) > 0 {
					runCommand(action.Value, ghaction.EnVars)
				}
			}
		}
	}
	return nil
}

func runCommand(command string, enVars []string) {
	envStr := strings.Join(enVars, "; ")
	if len(enVars) > 0 {
		envStr += ";"
	}
	cmd, err := exec.Command("sh", "-c", envStr+command).Output()
	if err != nil {
		log.Printf("Err: %s", err.Error())
		return
	}
	if *appDebug {
		log.Println(string(cmd))
	}
}
