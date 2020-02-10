package main

import (
	"fmt"
	"log"
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
