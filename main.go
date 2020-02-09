package main

import (
	"os"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	//Global flags
	app           = kingpin.New("whsub", "A WebHook subscriber")
	appDebug      = app.Flag("debug", "Enable debug mode.").Short('d').Bool()
	appConfigFile = app.Flag("config", "the configuration file for the subscriber").Short('c').File()

	//Server chlid command
	serverCmd        = app.Command("server", "Commands for the WH subscriber server")
	serverCmdStart   = serverCmd.Command("start", "Start the server")
	serverCmdVersion = serverCmd.Flag("version", "Show the version of the version of the server").Short('v').Bool()

	//Subsrcibe child command
	subscribeWh    = app.Command("subscribe", "Subscribe to a webhook").Alias("sub")
	subscribeWhID  = subscribeWh.Arg("whid", "Which webhook you want to subscribe").Required().String()
	subscribeWhURL = subscribeWh.Arg("url", "The URL to receive the notifications").Required().Envar("WH_URL").String()
)

func main() {
	app.HelpFlag.Short('h')
	if checkServerVersion() {
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

func checkServerVersion() bool {
	args := os.Args
	if len(args) == 3 && args[1] == serverCmd.FullCommand() && (args[2] == "--version" || args[2] == "-v") {
		printServerVersion()
		return true
	}
	return false
}
