package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	gaw "github.com/JojiiOfficial/GoAw"
	dbhelper "github.com/JojiiOfficial/GoDBHelper"
	"github.com/fatih/color"
)

//Actions the available actions
var Actions = map[string]int8{
	"github": 3, "gitlab": 1, "docker": 2, "script": 0,
}

func getWhIDFromHumanInput(db *dbhelper.DBhelper, input string) (int64, error) {
	whID := int64(-1)
	if len(input) > 0 && strings.Contains(input, "-") {
		realID := strings.Trim(strings.Split(input, "-")[1], " ")
		if len(realID) > 0 {
			var err error
			whID, err = getSubscriptionID(db, realID)
			if err != nil {
				whID = -1
				return whID, errors.New("no wh found")
			}
		}
	} else if !strings.Contains(input, "-") {
		return -1, errors.New("no wh found")
	}
	return whID, nil
}

//AddAction adds a new action
func AddAction(db *dbhelper.DBhelper, actionType, actionName, webhookName, actionFileDir string) {
	mode := Actions[actionType]

	hasAction, err := hasAction(db, actionName)
	if err != nil {
		log.Fatalln(err.Error())
		return
	}
	if hasAction {
		fmt.Printf("There is already an action with the name '%s'", actionName)
		return
	}

	scriptPathAbs, exists := gaw.DirAbs(actionFileDir)
	if !exists {
		log.Fatalf("Path '%s' does not exists", scriptPathAbs)
		return
	}

	var whID int64
	if len(webhookName) > 0 {
		whID, err = getWhIDFromHumanInput(db, webhookName)
		if err != nil {
			fmt.Printf(color.HiYellowString("Warning")+" webhook '%s' not found\n", webhookName)
		}
	}

	ending := ".yml"
	if mode == 0 {
		ending = ".sh"
	}
	file := path.Join(scriptPathAbs, actionName+ending)

	action := Action{
		Mode:           mode,
		File:           file,
		SubscriptionID: whID,
		Name:           actionName,
	}

	err = action.insert(db)
	if err != nil {
		log.Println(err.Error())
	}

	switch mode {
	case 0:
		{
			f, err := os.Create(file)
			if err != nil {
				log.Fatalln(err.Error())
				return
			}
			f.WriteString("#!/bin/bash\n")
			f.Close()
		}
	case 3:
		{
			if err := createDefaultGithubFile(file); err != nil {
				log.Fatalln(err.Error())
				return
			}
		}
	}

	fmt.Printf("Created action %s %s\n", actionName, color.HiGreenString("successfully"))
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
		mode := mapKeyByValue(action.Mode, Actions)
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
	err = updateActionWebhook(db, aID, whID)
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
		iMode, ok := Actions[newMode]
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

		if !gaw.FileExists(newFile) {
			y, i := gaw.ConfirmInput(color.HiYellowString("Warning: ")+"file does'n exist! "+"Continue anyway? [y/n]> ", bufio.NewReader(os.Stdin))
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

//Run an action
func (action *Action) Run(hookFile string) {
	if action.Mode == 0 {
		b, err := exec.Command(action.File).Output()
		if err != nil {
			log.Printf("Error executing action '%s': %s\n", action.Name, err.Error())
		} else if *appDebug {
			log.Printf("Output from %s:\n%s\n", action.Name, string(b))
		}
	}
}
