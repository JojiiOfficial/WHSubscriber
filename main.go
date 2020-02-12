package main

import (
	"log"
	"os"
	"strconv"

	godbhelper "github.com/JojiiOfficial/GoDBHelper"

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
	serverCmd      = app.Command("server", "Commands for the WH subscriber server")
	serverCmdStart = serverCmd.Command("start", "Start the server")

	//Subscriptions
	subscribtionsCmd = app.Command("subscriptions", "Lists your subscriptions").FullCommand()

	//Subsciption
	subscriptionCmd = app.Command("subscription", "Subscription command")
	//Subsrcibe child command
	subscribeAddWh         = subscriptionCmd.Command("add", "Subscribe to a webhook")
	subscribeWhID          = subscribeAddWh.Arg("webhookID", "Which webhook you want to subscribe").Required().String()
	subscribeWhCallbackURL = subscribeAddWh.Arg("url", "The callback URL to receive the notifications").Envar(getEnVar(EnVarReceiveURL)).String()
	subscribeWhScript      = subscribeAddWh.Flag("script", "The script to run on a webhook call").Short('s').String()

	subscribeImport = subscriptionCmd.Command("import", "Imports a subscription")

	//Config child command
	configCmd           = app.Command("config", "Commands for the config file")
	configCmdCreate     = configCmd.Command("create", "Create config file")
	configCmdCreateName = configCmdCreate.Arg("name", "Config filename").Required().String()

	//Actions
	actionsCmd = app.Command("actions", "List you actions").FullCommand()

	//Action commands
	actionCmd = app.Command("action", "Configure your actions for wehbooks")
	//Action add
	actionCmdAdd           = actionCmd.Command("add", "Adds an action for a webhook")
	actionCmdAddActionType = actionCmdAdd.Flag("action", "The kind of action you want to add").HintAction(hintAvailableActions).Default("script").String()
	actionCmdAddName       = actionCmdAdd.Flag("name", "The name of the action. To make it recycleable").HintAction(hintRandomNames).Default(getRandomName()).String()
	actionCmdAddWebhook    = actionCmdAdd.Flag("webhook", "The webhook to add the action to").HintAction(hintSubscriptions).String()
	actionCmdAddPath       = actionCmdAdd.Arg("dir", "The dir where the action-file should be created").HintAction(hintListCurrDir).Required().String()
	//Action setWebhook
	actionCmdSetWh        = actionCmd.Command("setwebhook", "Sets/Changes the webhook for an action")
	actionCmdSetWhAction  = actionCmdSetWh.Arg("action", "The action to change the webhook for").HintAction(hintListActions).Required().String()
	actionCmdSetWhWebhook = actionCmdSetWh.Arg("webhook", "The new webhook").HintAction(hintSubscriptions).Required().String()
	//Action delete
	actionCmdDelete     = actionCmd.Command("delete", "Deletes an action from a webhook").Alias("rm")
	actionCmdDeleteName = actionCmdDelete.Arg("name", "The name of the action").HintAction(hintListActions).Required().Strings()

	//Sources
	sourceCmd = app.Command("source", "Source command")
	//Infos for source
	sourceCmdInfos   = sourceCmd.Command("info", "Shows information for source")
	sourceCmdInfosID = sourceCmdInfos.Arg("sourceID", "The ID of the source to display informations for").String()
	//Create source
	sourceCmdCreate = sourceCmd.Command("create", "Create a new source")
	//Delete source
	sourceCmdDelete   = sourceCmd.Command("delete", "Delete a source")
	sourceCmdDeleteID = sourceCmdDelete.Arg("sourceID", "The ID of the source to delete").String()
)

func main() {
	app.HelpFlag.Short('h')
	app.Version(strconv.FormatFloat(SubscriberVersion, 'f', 2, 32))
	if checkVersionCommand() {
		return
	}

	//parsing the args
	parsed := kingpin.MustParse(app.Parse(os.Args[1:]))

	var (
		database = *appDatabase
		config   *ConfigStruct
		db       *godbhelper.DBhelper
	)

	if parsed != configCmdCreate.FullCommand() {
		//Return on error
		var shouldExit bool
		config, shouldExit = InitConfig(*appCfgFile, false)
		if shouldExit {
			return
		}
		if !config.Check() {
			if *appDebug {
				log.Println("Exiting")
			}
			return
		}
		var err error
		db, err = connectDB(database)
		if err != nil {
			log.Fatalln(err.Error())
			return
		}
	}

	switch parsed {
	//Server --------------------
	case serverCmdStart.FullCommand():
		{
			//whsub server start
			if config.CheckServer() {
				StartReceiverServer(config, db, *appDebug)
			}
		}

	//Subscriptions --------------------
	case subscribeAddWh.FullCommand():
		{
			//whsub subscription add
			Subscribe(db, config, *subscribeWhCallbackURL, *subscribeWhID)
		}
	case subscribtionsCmd:
		{
			//whsub subscriptions
			ViewSubscriptions(db)
		}

	//Actions --------------------
	case actionCmdAdd.FullCommand():
		{
			//whsub action add
			AddAction(db, *actionCmdAddActionType, *actionCmdAddName, *actionCmdAddWebhook, *actionCmdAddPath)
		}
	case actionCmdDelete.FullCommand():
		{
			//whsub action delete
			DeleteAction(db, *actionCmdDeleteName)
		}
	case actionsCmd:
		{
			//whsub actions
			ViewActions(db)
		}
	case actionCmdSetWh.FullCommand():
		{
			//whsub	action setwebhook
			ActionSetWebhook(db, *actionCmdSetWhWebhook, *actionCmdSetWhAction)
		}

	//Config --------------------
	case configCmdCreate.FullCommand():
		{
			//whsub config create
			InitConfig(*configCmdCreateName, true)
		}

	//Source --------------------
	case sourceCmdCreate.FullCommand():
		{
			//whsub source create
			CreateSource()
		}
	case sourceCmdDelete.FullCommand():
		{
			//whsub source delete
			DeleteSource()
		}
	case sourceCmdInfos.FullCommand():
		{
			//whsub source info
			SourceInfo()
		}
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
