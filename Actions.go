package main

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/fatih/color"
)

func getWhIDFromHumanInput(input string) (int64, error) {
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

func addAction() {
	mode := int8(0)
	if *actionCmdAddFMode == "action" {
		mode = 1
	}

	scriptPath := (*actionCmdAddAFile)
	scriptFileAbs, exists := fileFullPath(scriptPath)
	if !exists {
		log.Fatalf("File '%s' does not exists", scriptPath)
		return
	}

	whID, err := getWhIDFromHumanInput(*actionCmdAddWebhook)
	if err != nil {
		fmt.Println(color.HiYellowString("Warning"), "webhook-subscription not found")
	}

	action := Action{
		Mode:   mode,
		File:   scriptFileAbs,
		HookID: whID,
		Name:   *actionCmdAddName,
	}

	err = action.insert(db)
	if err != nil {
		log.Println(err.Error())
	}
	fmt.Printf("Created action %s %s\n", *actionCmdAddName, color.HiGreenString("successfully"))
}

func printActionList() {
	actions, err := getActions(db)
	if err != nil {
		log.Fatalln(err.Error())
		return
	}
	headingColor := color.New(color.FgHiGreen, color.Underline, color.Bold)
	headingColor.Println("ID\tName\t\tWebhook\t\t\tMode\tFile")
	for _, action := range actions {
		mode := "script"
		if action.Mode == 1 {
			mode = "action"
		}
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

func delAction() {
	actionIDs := (*actionCmdDeleteAID)
	for _, actionID := range actionIDs {
		has, err := hasAction(db, actionID)
		if err != nil {
			fmt.Println("An error occured:", err.Error())
			return
		}
		if !has {
			fmt.Printf("Action '%s' does %s\n", actionID, color.RedString("not exists"))
			return
		}
		err = deleteActionByID(db, actionID)
		if err == nil {
			fmt.Printf("Action '%s' deleted %s\n", actionID, color.HiGreenString("successful"))
		} else {
			fmt.Println("Error deleting action:", err.Error())
		}
	}
}

func actionSetWebhook() {
	whs := *actionCmdSetWhWebhook
	var whID int64
	var err error
	if whs != "na-" {
		whID, err = getWhIDFromHumanInput(*actionCmdSetWhWebhook)
		if err != nil {
			fmt.Println("Error webhook-subscription", color.HiRedString("not found"))
			return
		}
	} else {
		whID = -1
	}
	aID, err := getActionFromName(db, *actionCmdSetWhAction)
	if err != nil {
		fmt.Println("Error action", color.HiRedString("not found"))
		return
	}
	err = updateActionWebhook(db, aID, whID)
	if err != nil {
		log.Fatalln(err.Error())
		return
	}
	fmt.Printf("Action %s updated %s", *actionCmdSetWhAction, color.HiGreenString("successfully"))
}
