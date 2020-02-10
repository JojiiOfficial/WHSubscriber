package main

import (
	"log"
	"os"

	dbhelper "github.com/JojiiOfficial/GoDBHelper"
	"github.com/gobuffalo/packr"
	_ "github.com/mattn/go-sqlite3"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	//Global flags
	app         = kingpin.New("whsub", "A WebHook subscriber")
	appDebug    = app.Flag("debug", "Enable debug mode.").Short('d').Bool()
	appNoColor  = app.Flag("no-color", "Disable colors").Bool()
	appDatabase = app.Flag("database", "Path to the database to use").Default(getDefaultDBFile()).Envar(getEnVar(EnVarDatabaseFile)).String()
	appCfgFile  = app.
			Flag("config", "the configuration file for the subscriber").
			Envar(getEnVar(EnVarConfigFile)).
			Short('c').String()

	//Server chlid command
	serverCmd         = app.Command("server", "Commands for the WH subscriber server")
	serverCmdCStart   = serverCmd.Command("start", "Start the server")
	serverCmdFVersion = serverCmd.Flag("version", "Show the version of the server").Bool()

	//Subsrcibe child command
	subscribeWh             = app.Command("subscribe", "Subscribe to a webhook")
	subscribeWhAID          = subscribeWh.Arg("webhookID", "Which webhook you want to subscribe").Required().String()
	subscribeWhACallbackURL = subscribeWh.Arg("url", "The callback URL to receive the notifications").Envar(getEnVar(EnVarReceiveURL)).String()
	subscribeWhFScript      = subscribeWh.Flag("script", "The script to run on a webhook call").Short('s').String()

	//Config child command
	configCmd            = app.Command("config", "Commands for the config file")
	configCmdCCreate     = configCmd.Command("create", "Create config file")
	configCmdACreateName = configCmdCCreate.Arg("name", "Config filename").Required().String()

	//Action commands
	actionCmd = app.Command("actions", "Configure your actions for wehbooks")
	//Action add
	actionCmdCAdd       = actionCmd.Command("add", "Adds an action for a webhook")
	actionCmdAddFMode   = actionCmdCAdd.Flag("mode", "The kind of action you want to add (script / action)").HintOptions("script", "action").Default("script").String()
	actionCmdAddName    = actionCmdCAdd.Flag("name", "The name of the action. To make it recycleable").HintAction(hintRandomNames).Default(getRandomName()).String()
	actionCmdAddWebhook = actionCmdCAdd.Arg("webhook", "The webhook to add the action to").HintAction(hintSubscriptions).Required().String()
	actionCmdAddAFile   = actionCmdCAdd.Arg("file", "the file of the action (a script or action file)").HintAction(hintListCurrDir).Required().String()
	//Action list
	actionCmdCList = actionCmd.Command("list", "lists the actions")
	//Action delete
	actionCmdCDelete   = actionCmd.Command("delete", "Deletes an action from a webhook")
	actionCmdDeleteAID = actionCmdCDelete.Arg("id", "The name of the script").Required().Int()
)

var (
	db         *dbhelper.DBhelper
	database   string
	appDataBox *packr.Box
)

func main() {
	app.HelpFlag.Short('h')
	if checkVersionCommand() {
		return
	}

	//parsing the args
	parsed := kingpin.MustParse(app.Parse(os.Args[1:]))
	database = *appDatabase

	if parsed != configCmdCCreate.FullCommand() {
		//Return on error
		if InitConfig(*appCfgFile, false) {
			return
		}
		if !config.Check() {
			if *appDebug {
				log.Println("Exiting")
			}
			return
		}
		if err := connectDB(); err != nil {
			log.Fatalln(err.Error())
			return
		}
	}

	//Runnig the correct child command
	switch parsed {
	case serverCmdCStart.FullCommand():
		runWHReceiverServer()
	case subscribeWh.FullCommand():
		subscribe()
	case configCmdCCreate.FullCommand():
		InitConfig(*configCmdACreateName, true)
	case actionCmdCAdd.FullCommand():
		addAction()
	case actionCmdCList.FullCommand():
		printActionList()
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
