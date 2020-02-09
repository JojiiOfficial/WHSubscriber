package main

import (
	"fmt"
)

//Env vars
const (
	//EnVarPrefix prefix of all used env vars
	EnVarPrefix = "WHS"

	EnVarReceiveURL = "URL"
	EnVarConfigFile = "CONFIG"
)

func getEnVar(name string) string {
	return fmt.Sprintf("%s_%s", EnVarPrefix, name)
}
