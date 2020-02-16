package main

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/JojiiOfficial/configor"
)

//ActionFileStruct the struct for action files
type ActionFileStruct struct {
	Source  SourceActionItem
	Shell   ShellActionItem
	Actions []string
}

//SourceActionItem item for an action
type SourceActionItem struct {
	Trigger string
	Filter  map[string]string
}

//ShellActionItem action item for shell options
type ShellActionItem struct {
	User   string
	EnVars []string
}

func createDefaultActionFile(file string) error {
	_, err := os.Stat(file)
	//Return if already exists
	if err == nil {
		return nil
	}

	//Get Path
	pwd, _ := filepath.Split(file)
	if strings.HasSuffix(pwd, "/") {
		pwd = pwd[:len(pwd)-1]
	}

	username := getUsername()

	actionFileStruct := ActionFileStruct{
		Source: SourceActionItem{
			Trigger: "push",
			Filter:  map[string]string{"branch": "master"},
		},

		Shell: ShellActionItem{
			User: username,
			EnVars: []string{
				"PATH=/bin:/sbin:/usr/bin:/usr/local/sbin:/usr/local/bin",
				"ACTION_PWD=" + pwd,
			},
		},

		Actions: []string{
			"/some/script/Here",
		},
	}

	_, err = configor.SetupConfig(&actionFileStruct, file, configor.NoChange)
	return err
}

//LoadActionFile loads the action from a file
func LoadActionFile(file string) (*ActionFileStruct, error) {
	action := ActionFileStruct{}
	err := configor.Load(&action, file)
	if err != nil {
		return nil, err
	}
	return &action, nil
}

//Run runs the action-file
func (action *ActionFileStruct) Run(payloadFile string, saction *Action) error {
	if len(action.Actions) == 0 {
		return errors.New("No action defined")
	}

	for _, actionCmd := range action.Actions {
		if len(actionCmd) == 0 {
			continue
		}

		username := getUsername(action.Shell.User)

		runCommand(actionCmd, username, saction, action.Shell.EnVars)
	}

	return nil
}

func runCommand(command, username string, action *Action, enVars []string) {
	command = replaceRelativePath(command, action)
	envStr := formatBashEnVars(enVars)

	var execCommand string
	var args []string

	//Use 'su' if running as root to allow switching users
	if os.Getuid() == 0 {
		execCommand = "su"
		args = []string{username, "-c", envStr + command}
	} else {
		execCommand = "sh"
		args = []string{"-c", envStr + command}
	}
	cmd, err := exec.Command(execCommand, args...).Output()

	if err != nil {
		log.Printf("Err: %s", err.Error())
	} else if *appDebug {
		log.Println("Output from '" + action.Name + "':\n" + string(cmd))
	}
}

//Use action-folder for relative file instead of binary-relative path
func replaceRelativePath(file string, action *Action) string {
	if strings.HasPrefix(file, "./") {
		pah, _ := filepath.Split(action.File)
		file = path.Join(pah, file)
	}
	return file
}
