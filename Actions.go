package main

import (
	"fmt"
	"log"

	"github.com/fatih/color"
)

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

	action := Action{
		Mode:   mode,
		File:   scriptFileAbs,
		HookID: "0",
		Name:   *actionCmdAddName,
	}

	err := action.insert(db)
	if err != nil {
		log.Println(err.Error())
	}
	fmt.Println(action.ID)
}

func printActionList() {
	actions, err := getActions(db)
	if err != nil {
		log.Fatalln(err.Error())
		return
	}
	headingColor := color.New(color.FgHiGreen, color.Underline, color.Bold)
	headingColor.Println("ID\tName\t\tWebhook\t\tMode\tFile")
	for _, action := range actions {
		mode := "script"
		if action.Mode == 1 {
			mode = "action"
		}
		name := action.Name
		if len(name) < 8 {
			name += "\t"
		}
		fmt.Printf("%d\t%s\t%s\t\t%s\t%s\n", action.ID, name, action.HookID, mode, action.File)
	}
	if len(actions) == 0 {
		fmt.Println("No action available")
	}
	fmt.Println()
}

func delAction() {
	actionID := (*actionCmdDeleteAID)
	has, err := hasAction(db, actionID)
	if err != nil {
		fmt.Println("An error occured:", err.Error())
		return
	}
	if !has {
		fmt.Println("Action does not exists")
		return
	}
	err = deleteActionByID(db, actionID)
	if err == nil {
		fmt.Println("Action deleted successfully")
	} else {
		fmt.Println("Error deleting action:", err.Error())
	}
}
