package main

import (
	"os"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	//Global flags
	app           = kingpin.New("whsub", "A WebHook subscriber")
	appDebug      = app.Flag("debug", "Enable debug mode.").Short('d').Bool()
	appConfigFile = app.Flag("config", "the configuration file for the subscriber").Envar(getEnVar(EnVarConfigFile)).Short('c').File()

	//Server chlid command
	serverCmd        = app.Command("server", "Commands for the WH subscriber server")
	serverCmdStart   = serverCmd.Command("start", "Start the server")
	serverCmdVersion = serverCmd.Flag("version", "Show the version of the server").Bool()

	//Subsrcibe child command
	subscribeWh    = app.Command("subscribe", "Subscribe to a webhook")
	subscribeWhID  = subscribeWh.Arg("whid", "Which webhook you want to subscribe").Required().String()
	subscribeWhURL = subscribeWh.Arg("url", "The URL to receive the notifications").Envar(getEnVar(EnVarReceiveURL)).String()
)

func main() {
	app.HelpFlag.Short('h')
	if checkVersionCommand() {
		return
	}

	//Runnig the correct child command
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case serverCmdStart.FullCommand():
		runWHReceiverServer()
	case subscribeWh.FullCommand():
		subscribe()
	}
}

func checkVersionCommand() bool {
	args := os.Args
	if len(args) == 3 && (args[2] == "--version" || args[2] == "-v") {
		switch args[1] {
		case serverCmd.FullCommand():
			printServerVersion()
			return true
		case subscribeWh.FullCommand(), "subscriber", "sub":
			printSubscriberVersion()
			return true
		}
	}
	return false
}
