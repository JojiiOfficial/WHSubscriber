package main

//ActionItem an item in an action file
type ActionItem struct {
	Type  ActionType
	Value string
}

//ActionType the type of action
type ActionType string

const (
	//CommandActionItem a single command to run
	CommandActionItem ActionType = "command"
	//ScriptActionItem a script to run
	ScriptActionItem ActionType = "script"
)
