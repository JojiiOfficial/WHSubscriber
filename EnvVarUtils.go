package main

import (
	"fmt"
)

//Env vars
const (
	//EnvVarPrefix prefix of all used env vars
	EnvVarPrefix = "WHS"

	EnvVarReceiveURL = "URL"
	EnvVarConfigFile = "CONFIG"
)

func getEnvVar(name string) string {
	return fmt.Sprintf("%s_%s", EnvVarPrefix, name)
}
