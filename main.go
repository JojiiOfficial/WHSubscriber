package main

import (
	"os"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	//Global flags
	app        = kingpin.New("whsub", "A WebHook subscriber")
	appDebug   = app.Flag("debug", "Enable debug mode.").Short('d').Bool()
	appCfgFile = app.
			Flag("config", "the configuration file for the subscriber").
			Envar(getEnVar(EnVarConfigFile)).
			Short('c').String()

	//Server chlid command
	serverCmd        = app.Command("server", "Commands for the WH subscriber server")
	serverCmdStart   = serverCmd.Command("start", "Start the server")
	serverCmdVersion = serverCmd.Flag("version", "Show the version of the server").Bool()

	//Subsrcibe child command
	subscribeWh            = app.Command("subscribe", "Subscribe to a webhook")
	subscribeWhID          = subscribeWh.Arg("whid", "Which webhook you want to subscribe").Required().String()
	subscribeWhCallbackURL = subscribeWh.Arg("url", "The callback URL to receive the notifications").Envar(getEnVar(EnVarReceiveURL)).String()
	subscribeWhScript      = subscribeWh.Flag("script", "The script to run on a webhook call").Short('s').String()

	//Config child command
	configCmd           = app.Command("config", "Commands for the config file")
	configCmdCreate     = configCmd.Command("create", "Create config file")
	configCmdCreateName = configCmdCreate.Arg("name", "Config filename").Required().String()
)

func main() {
	app.HelpFlag.Short('h')
	if checkVersionCommand() {
		return
	}

	//parsing the args
	parsed := kingpin.MustParse(app.Parse(os.Args[1:]))

	if parsed != configCmdCreate.FullCommand() {
		//Return on error
		if InitConfig(*appCfgFile, false) {
			return
		}
	}

	//Runnig the correct child command
	switch parsed {
	case serverCmdStart.FullCommand():
		runWHReceiverServer()
	case subscribeWh.FullCommand():
		subscribe()
	case configCmdCreate.FullCommand():
		InitConfig(*configCmdCreateName, true)
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
