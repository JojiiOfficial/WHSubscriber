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
	appDebug    = app.Flag("debug", "Enable debug mode").Short('d').Bool()
	appNoColor  = app.Flag("no-color", "Disable colors").Envar(getEnVar(EnVarNoColor)).Bool()
	appDatabase = app.Flag("database", "Path to the database to use").Default(getDefaultDBFile()).Envar(getEnVar(EnVarDatabaseFile)).String()
	appCfgFile  = app.
			Flag("config", "the configuration file for the subscriber").
			Envar(getEnVar(EnVarConfigFile)).
			Short('c').String()

	//Server chlid command
	serverCmd         = app.Command("server", "Commands for the WH subscriber server")
	serverCmdCStart   = serverCmd.Command("start", "Start the server")
	serverCmdFVersion = serverCmd.Flag("version", "Show the version of the server").Bool()

	//Subscriptions
	subscribtions = app.Command("subscriptions", "Lists your subscriptions").FullCommand()

	//Subsciption
	subscription = app.Command("subscription", "Subscription command")
	//Subsrcibe child command
	subscribeAddWh          = subscription.Command("add", "Subscribe to a webhook")
	subscribeWhAID          = subscribeAddWh.Arg("webhookID", "Which webhook you want to subscribe").Required().String()
	subscribeWhACallbackURL = subscribeAddWh.Arg("url", "The callback URL to receive the notifications").Envar(getEnVar(EnVarReceiveURL)).String()
	subscribeWhFScript      = subscribeAddWh.Flag("script", "The script to run on a webhook call").Short('s').String()

	subscribeImport = subscription.Command("import", "Imports a subscription")

	//Config child command
	configCmd            = app.Command("config", "Commands for the config file")
	configCmdCCreate     = configCmd.Command("create", "Create config file")
	configCmdACreateName = configCmdCCreate.Arg("name", "Config filename").Required().String()

	//Actions
	actionsCmd = app.Command("actions", "Configure your actions for wehbooks").FullCommand()

	//Action commands
	actionCmd = app.Command("action", "Configure your actions for wehbooks")
	//Action add
	actionCmdCAdd       = actionCmd.Command("add", "Adds an action for a webhook")
	actionCmdAddFAction = actionCmdCAdd.Flag("action", "The kind of action you want to add").HintAction(hintAvailableActions).Default("script").String()
	actionCmdAddName    = actionCmdCAdd.Flag("name", "The name of the action. To make it recycleable").HintAction(hintRandomNames).Default(getRandomName()).String()
	actionCmdAddWebhook = actionCmdCAdd.Flag("webhook", "The webhook to add the action to").HintAction(hintSubscriptions).String()
	actionCmdAddAFile   = actionCmdCAdd.Arg("file", "the file of the action (a script or action file)").HintAction(hintListCurrDir).Required().String()
	//Action setWebhook
	actionCmdCSetWh       = actionCmd.Command("setwebhook", "Sets/Changes the webhook for an action")
	actionCmdSetWhAction  = actionCmdCSetWh.Arg("action", "The action to change the webhook for").HintAction(hintListActions).Required().String()
	actionCmdSetWhWebhook = actionCmdCSetWh.Arg("webhook", "The new webhook").HintAction(hintSubscriptions).Required().String()
	//Action delete
	actionCmdCDelete   = actionCmd.Command("delete", "Deletes an action from a webhook").Alias("rm")
	actionCmdDeleteAID = actionCmdCDelete.Arg("name", "The name of the action").HintAction(hintListActions).Required().Strings()
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
	case subscribeAddWh.FullCommand():
		subscribe()
	case configCmdCCreate.FullCommand():
		InitConfig(*configCmdACreateName, true)
	case actionCmdCAdd.FullCommand():
		addAction()
	case actionCmdCDelete.FullCommand():
		delAction()
	case actionsCmd:
		printActionList()
	case subscribtions:
		printSubsciptionList()
	case actionCmdCSetWh.FullCommand():
		actionSetWebhook()
	}
}

func checkVersionCommand() bool {
	args := os.Args
	if len(args) == 3 && (args[2] == "--version" || args[2] == "-v") {
		switch args[1] {
		case serverCmd.FullCommand():
			printServerVersion()
			return true
		case subscribeAddWh.FullCommand(), "subscriber", "sub":
			printSubscriberVersion()
			return true
		}
	}
	return false
}
