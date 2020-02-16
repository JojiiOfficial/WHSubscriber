package main

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path"
	"path/filepath"
	"strings"

	"github.com/JojiiOfficial/configor"
)

//GitActionStruct the struct for github webhooks
type GitActionStruct struct {
	Git     GitActionItem
	Shell   ShellActionItem
	Actions []string
}

//GitRemoteServer remote server for git (github/gitlab/etc...)
type GitRemoteServer string

const (
	//Github github.com
	Github GitRemoteServer = "github"
	//Gitlab gitlab.com
	Gitlab GitRemoteServer = "gitlab"
)

//GitActionItem item for git action
type GitActionItem struct {
	RemoteServer GitRemoteServer
	Trigger      string
	Filter       map[string]string
}

//ShellActionItem action item for shell options
type ShellActionItem struct {
	User   string
	EnVars []string
}

func createDefaultGitFile(file string, gitServer GitRemoteServer) error {
	_, err := os.Stat(file)
	if err == nil {
		return nil
	}
	pwd, _ := filepath.Split(file)
	if strings.HasSuffix(pwd, "/") {
		pwd = pwd[:len(pwd)-1]
	}
	var username string
	user, err := user.Current()
	if err == nil {
		username = user.Username
	}
	gitActionStruct := GitActionStruct{
		Git: GitActionItem{
			RemoteServer: gitServer,
			Trigger:      "push",
			Filter:       map[string]string{"branch": "master"},
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
	_, err = configor.SetupConfig(&gitActionStruct, file, configor.NoChange)
	return err
}

//LoadGitAction loads the action from a file
func LoadGitAction(file string) (*GitActionStruct, error) {
	action := GitActionStruct{}
	err := configor.Load(&action, file)
	if err != nil {
		return nil, err
	}
	return &action, nil
}

//Run runs the github action
func (gitAction *GitActionStruct) Run(payloadFile string, saction *Action) error {
	if len(gitAction.Actions) == 0 {
		return errors.New("no action defined")
	}

	//TODO Parse the incoming action

	for _, action := range gitAction.Actions {
		if len(action) == 0 {
			continue

		}
		username := gitAction.Shell.User
		if len(username) == 0 {
			user, err := user.Current()
			if err != nil {
				log.Println(err.Error())
				continue
			}
			username = user.Username
		}
		runCommand(action, username, saction, gitAction.Shell.EnVars)

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

func runCommand(command, username string, action *Action, enVars []string) {
	command = replaceRelativePath(command, action)
	envStr := formatenvvars(enVars)

	var execCommand string
	var args []string
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
