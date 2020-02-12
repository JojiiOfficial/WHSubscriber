package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	godbhelper "github.com/JojiiOfficial/GoDBHelper"
	"github.com/fatih/color"
)

//Actions the available actions
var Actions = map[string]int8{
	"github": 3, "gitlab": 1, "docker": 2, "script": 0,
}

func getWhIDFromHumanInput(db *godbhelper.DBhelper, input string) (int64, error) {
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
func AddAction(db *godbhelper.DBhelper, actionType, actionName, webhookName, actionFileDir string) {
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

	scriptPathAbs, exists := dirAbs(actionFileDir)
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
func ViewActions(db *godbhelper.DBhelper) {
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
func DeleteAction(db *godbhelper.DBhelper, actionIDs []string) {
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
func ActionSetWebhook(db *godbhelper.DBhelper, webhookName, actionName string) {
	var whID int64
	var err error
	if webhookName != "na-" {
		whID, err = getWhIDFromHumanInput(db, webhookName)
		if err != nil {
			fmt.Println("Error webhook-subscription", color.HiRedString("not found"))
			return
		}
	} else {
		whID = -1
	}
	aID, err := getActionFromName(db, actionName)
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
