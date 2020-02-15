package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	gaw "github.com/JojiiOfficial/GoAw"
	dbhelper "github.com/JojiiOfficial/GoDBHelper"
	"github.com/fatih/color"
)

//Modes the available actions
var Modes = map[string]uint8{
	"github": 3, "gitlab": 1, "docker": 2, "script": 0,
}

func getWhIDFromHumanInput(db *dbhelper.DBhelper, input string) (int64, error) {
	whID := int64(-1)
	realID := ""
	if len(input) > 0 && strings.Contains(input, "-") {
		realID = strings.Trim(strings.Split(input, "-")[1], " ")
	} else if !strings.Contains(input, "-") {
		realID = input
	}

	if len(realID) > 0 {
		var err error
		whID, err = getSubscriptionID(db, realID)
		if err != nil {
			whID = -1
			return whID, errors.New("no wh found")
		}
	} else {
		return -1, errors.New("no wh found")
	}
	return whID, nil
}

//AddAction adds a new action
func AddAction(db *dbhelper.DBhelper, actionType, actionName, webhookName, actionFile string, createFile bool) {
	mode := Modes[actionType]

	hasAction, err := hasAction(db, actionName)
	if err != nil {
		log.Fatalln(err.Error())
		return
	}
	if hasAction {
		fmt.Printf("There is already an action with the name '%s'", actionName)
		return
	}

	var whID int64
	if len(webhookName) > 0 {
		whID, err = getWhIDFromHumanInput(db, webhookName)
		if err != nil {
			fmt.Printf(color.HiYellowString("Warning")+" webhook '%s' not found\n", webhookName)
		}
	}

	newFileString := gaw.FromString(actionFile)
	absfile, _ := filepath.Abs(actionFile)

	if !gaw.FileExists(newFileString.ToString()) && !(*appYes) && !createFile {
		y, i := gaw.ConfirmInput(color.HiYellowString("Warning: ")+"file or directory does'n exist! Continue anyway? [y/n]> ", bufio.NewReader(os.Stdin))
		if i == -1 || !y {
			fmt.Println("Abort")
			return
		}
	}

	if (newFileString.EndsWith(".sh") && mode != 0) || ((newFileString.EndsWith(".yml") || newFileString.EndsWith(".yaml")) && mode == 0) {
		fmt.Println(color.HiRedString("Err:"), " If you file is a script it needs to end with '.sh'!")
		return
	}

	action := Action{
		Mode:           mode,
		File:           absfile,
		SubscriptionID: whID,
		Name:           actionName,
	}

	err = action.insert(db)
	if err != nil {
		log.Println(err.Error())
	}

	fmt.Printf("Created action %s %s\n", actionName, color.HiGreenString("successfully"))

	if createFile {
		ActionCreateFile(db, &action)
	}
}

//ViewActions prints actions
func ViewActions(db *dbhelper.DBhelper) {
	actions, err := getActions(db)
	if err != nil {
		log.Fatalln(err.Error())
		return
	}
	headingColor := color.New(color.FgHiGreen, color.Underline, color.Bold)
	headingColor.Println("ID\tName\t\tWebhook\t\t\tMode\tFile")
	for _, action := range actions {
		mode := mapKeyByValue(action.Mode, Modes)
		name := action.Name
		if len(name) < 8 {
			name += "\t"
		}
		if len(action.HookName) < 8 {
			action.HookName += "\t"
		}
		fmt.Printf("%d\t%s\t%s\t\t%s\t%s\n", action.ID, name, action.HookName, mode, action.File)
	}
	if len(actions) == 0 {
		fmt.Println("No action available")
	}
	fmt.Println()
}

//DeleteAction deletes an action
func DeleteAction(db *dbhelper.DBhelper, actionIDs []string) {
	for _, actionID := range actionIDs {
		has, err := hasAction(db, actionID)
		if err != nil {
			fmt.Println("An error occured:", err.Error())
			return
		}
		if !has {
			fmt.Printf("Action '%s' does %s\n", actionID, color.RedString("not exist")+". Skipping...")
			continue
		}
		action, e := getActionFromName(db, actionID)
		if e != nil {
			fmt.Println(e.Error())
			continue
		}

		var delete bool
		if gaw.FileExists(action.File) && !(*appYes) {
			y, i := gaw.ConfirmInput("Delete action file for "+action.Name+"? [y/n]> ", bufio.NewReader(os.Stdin))
			if i == -1 {
				return
			}
			delete = y
		}
		if delete {
			os.Remove(action.File)
			if *appDebug {
				fmt.Printf("File %s deleted\n", action.File)
			}
		}
		err = deleteActionByID(db, actionID)
		if err == nil {
			fmt.Printf("Action '%s' deleted %s\n", actionID, color.HiGreenString("successful"))
		} else {
			fmt.Println(color.RedString("Err"), "deleting action:", err.Error())
		}
	}
}

//ActionSetWebhook sets webhook for an action
func ActionSetWebhook(db *dbhelper.DBhelper, webhookName, actionName string) {
	var whID int64
	var err error
	if webhookName != "na" {
		whID, err = getWhIDFromHumanInput(db, webhookName)
		if err != nil {
			fmt.Println("Error webhook-subscription", color.HiRedString("not found"))
			return
		}
	} else {
		whID = -1
	}
	aID, err := getActionIDFromName(db, actionName)
	if err != nil {
		fmt.Println("Error action", color.HiRedString("not found"))
		return
	}
	err = updateActionSource(db, aID, whID)
	if err != nil {
		log.Fatalln(err.Error())
		return
	}
	fmt.Printf("Action %s updated %s", actionName, color.HiGreenString("successfully"))
}

//ActionSetFile sets the actionfile for an action
func ActionSetFile(db *dbhelper.DBhelper, actionName, newMode, newFile string) {
	if len(newMode) == 0 && len(newFile) == 0 {
		fmt.Println("You need to set one of the flags:", color.HiRedString("--new-file"), "or", color.HiRedString("--new-mode"))
		return
	}
	action, err := getActionFromName(db, actionName)
	if err != nil {
		fmt.Println("Error action", color.HiRedString("not found"))
		return
	}

	if len(newMode) > 0 {
		iMode, ok := Modes[newMode]
		if !ok {
			fmt.Printf("Mode '%s' doesn't exist!\nSkipping\n", newMode)
		} else {
			err = updateActionMode(db, action.ID, iMode)
			if err != nil {
				log.Fatalln(err.Error())
				return
			}
			action.Mode = iMode
			fmt.Printf("Action '%s' %s updated to mode '%s'\n", actionName, color.HiGreenString("successfully"), color.HiGreenString(newMode))
		}
	}

	if len(newFile) > 0 {
		newFileString := gaw.FromString(newFile)
		absfile, _ := filepath.Abs(newFile)

		if !gaw.FileExists(newFile) && !(*appYes) {
			y, i := gaw.ConfirmInput(color.HiYellowString("Warning: ")+"file does'n exist! Continue anyway? [y/n]> ", bufio.NewReader(os.Stdin))
			if i == -1 || !y {
				fmt.Println("Abort")
				return
			}
		}

		if (newFileString.EndsWith(".sh") && action.Mode != 0) || ((newFileString.EndsWith(".yml") || newFileString.EndsWith(".yaml")) && action.Mode == 0) {
			fmt.Println(color.HiRedString("Err:"), " If you file is a script it needs to end with '.sh'!")
			return
		}

		err = updateActionFile(db, action.ID, absfile)
		if err != nil {
			log.Fatalln(err.Error())
			return
		}
		fmt.Printf("Action '%s' %s updated to file '%s'\n", actionName, color.HiGreenString("successfully"), color.HiGreenString(absfile))
	}
}

//ActionCreateFileFromName creates the file for an action by its name given
func ActionCreateFileFromName(db *dbhelper.DBhelper, actionName string) {
	action, err := getActionFromName(db, actionName)
	if err != nil {
		fmt.Println("Action", color.RedString("not found"))
		return
	}
	ActionCreateFile(db, action)
}

//ActionCreateFile creates the file for an action
func ActionCreateFile(db *dbhelper.DBhelper, action *Action) {
	if gaw.FileExists(action.File) && !(*appYes) {
		y, i := gaw.ConfirmInput(color.HiYellowString("Warning: ")+"file already exist! Overwrite? [y/n]> ", bufio.NewReader(os.Stdin))
		if i == -1 || !y {
			fmt.Println("Abort")
			return
		}
	}
	switch action.Mode {
	case 0:
		{
			//Bash script
			f, err := os.Create(action.File)
			if err != nil {
				log.Fatalln(err.Error())
				return
			}
			f.WriteString("#!/bin/bash\n")
			f.Close()
		}
	case 3, 1:
		{
			//Git(hub|lab)
			remote := Github
			if action.Mode == 1 {
				remote = Gitlab
			}
			if err := createDefaultGitFile(action.File, remote); err != nil {
				log.Fatalln(err.Error())
				return
			}
		}
	case 2:
		{
			//Docker
			//TODO implement create docker config
		}
	default:
		return
	}
	fmt.Println("Action file created", color.HiGreenString("sucessfully"))
}

//Run an action
func (action *Action) Run(hookFile string) {
	if !gaw.FileExists(action.File) {
		log.Println("Error: Action file", action.File, "doesn't exist!")
		return
	}

	switch action.Mode {
	case 0:
		{
			//Script
			b, err := exec.Command(action.File, hookFile).Output()
			if err != nil {
				log.Printf("Error executing action '%s': %s\n", action.Name, err.Error())
			} else if *appDebug {
				log.Printf("Output from %s:\n%s\n", action.Name, string(b))
			}
		}
	case 3, 1:
		{
			//Git(hub|lab)
			gitAction, err := LoadGitAction(action.File)
			if err != nil {
				log.Println(err.Error())
				return
			}
			gitAction.Run(hookFile, action)
		}
	case 2:
		{
			//Docker
			//TODO run docker action
		}
	}
}
