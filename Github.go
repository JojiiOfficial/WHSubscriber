package main

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"path/filepath"
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
	pwd, _ := filepath.Split(file)
	if strings.HasSuffix(pwd, "/") {
		pwd = pwd[:len(pwd)-1]
	}
	ghActionStruct := GithubActionStruct{
		Trigger: "push",
		Filter:  map[string]string{"branch": "master"},
		EnVars: []string{
			"PATH=/bin:/sbin:/usr/bin:/usr/local/sbin:/usr/local/bin",
			"ACTION_PWD=" + pwd,
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
func (ghaction *GithubActionStruct) Run(payloadFile, actionName string) error {
	if len(ghaction.Actions) == 0 {
		return errors.New("no action defined")
	}
	for _, action := range ghaction.Actions {
		if len(action.Value) == 0 {
			continue
		}
		switch action.Type {
		case ScriptActionItem:
			{
				runScript(action.Value, actionName, ghaction.EnVars)
			}
		case CommandActionItem:
			{
				runCommand(action.Value, actionName, ghaction.EnVars)
			}
		}
	}
	return nil
}

func formatenvvars(enVars []string) string {
	envStr := strings.Join(enVars, "; ")
	if len(enVars) > 0 {
		envStr += ";"
	}
	return envStr
}

func runCommand(command, actionName string, enVars []string) {
	envStr := formatenvvars(enVars)
	if *appDebug {
		log.Println("sh -c '" + envStr + command + "'")
	}

	cmd, err := exec.Command("sh", "-c", envStr+command).Output()
	if err != nil {
		log.Printf("Err: %s", err.Error())
		return
	}
	if *appDebug {
		log.Println("Output from '" + actionName + "':\n" + string(cmd))
	}
}

func runScript(file, actionName string, enVars []string) {

}
