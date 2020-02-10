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
	action := Action{
		Mode:   mode,
		File:   (*actionCmdAddAFile).Name(),
		HookID: "0",
		Name:   *actionCmdAddName,
	}

	err := action.insert(db)
	if err != nil {
		log.Println(err.Error())
	}
	fmt.Println(action.ID)
}
