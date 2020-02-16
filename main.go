package main

import (
	"log"
	"os"
	"strconv"

	dbhelper "github.com/JojiiOfficial/GoDBHelper"

	_ "github.com/mattn/go-sqlite3"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

const (
	//ServerVersion the version of the server
	ServerVersion = 0.2
	//SubscriberVersion version of the WebHook subscriber
	SubscriberVersion = 0.6
)

var (
	//Global command
	app = kingpin.New("whsub", "A WebHook subscriber")

	//Global flags
	appDebug    = app.Flag("debug", "Enable debug mode").Short('d').Bool()
	appNoColor  = app.Flag("no-color", "Disable colors").Envar(getEnVar(EnVarNoColor)).Bool()
	appYes      = app.Flag("yes", "Skips confirmations").Short('y').Envar(getEnVar(EnVarYes)).Bool()
	appDatabase = app.Flag("database", "Path to the database to use").Default(getDefaultDBFile()).Envar(getEnVar(EnVarDatabaseFile)).String()
	appCfgFile  = app.
			Flag("config", "the configuration file for the subscriber").
			Envar(getEnVar(EnVarConfigFile)).
			Short('c').String()

	//Server commands
	//Server start
	serverCmd      = app.Command("server", "Commands for the WH subscriber server")
	serverCmdStart = serverCmd.Command("start", "Start the server")

	//Subscription commands
	//Subscriptions
	subscriptionsCmd = app.Command("subscriptions", "Lists your subscriptions").FullCommand()
	//Subscription
	subscriptionCmd = app.Command("subscription", "Subscription command")
	//Subscribe child command
	subscribeAddWh         = subscriptionCmd.Command("add", "Subscribe to a webhook")
	subscribeWhID          = subscribeAddWh.Arg("webhookID", "Which webhook you want to subscribe").Required().String()
	subscribeWhCallbackURL = subscribeAddWh.Arg("url", "The callback URL to receive the notifications").Envar(getEnVar(EnVarReceiveURL)).String()
	subscribeWhScript      = subscribeAddWh.Flag("script", "The script to run on a webhook call").Short('s').String()
	//Subscribe import
	subscribeImport   = subscriptionCmd.Command("import", "Imports a subscription").Alias("load")
	subscribeImportID = subscribeImport.Arg("id", "The ID of the subscription to import").Required().String()
	//Subscription delete
	subscribeDelete   = subscriptionCmd.Command("unsubscribe", "Delete/unsubscribe a subscription").Alias("rm").Alias("delete")
	subscribeDeleteID = subscribeDelete.Arg("id", "The ID of the subscription to delete").HintAction(hintSubscriptionsNoNa).Required().String()

	//Config commands
	//Config create
	configCmd           = app.Command("config", "Commands for the config file")
	configCmdCreate     = configCmd.Command("create", "Create config file")
	configCmdCreateName = configCmdCreate.Arg("name", "Config filename").Required().String()

	//Action commands
	//Action
	actionCmd = app.Command("action", "Configure your actions for webhooks")
	//Actions
	actionsCmd = app.Command("actions", "List you actions").FullCommand()
	//Action add
	actionCmdAdd        = actionCmd.Command("add", "Adds an action for a webhook")
	actionCmdAddMode    = actionCmdAdd.Flag("mode", "The kind of action you want to add").HintAction(hintAvailableActions).String()
	actionCmdAddFile    = actionCmdAdd.Arg("file", "The action-file. Either a bash script or an action-configuration").HintAction(hintListCurDir).Required().String()
	actionCmdAddName    = actionCmdAdd.Flag("name", "The name of the action. To make it easier to use again").HintAction(hintRandomNames).Default(getRandomName()).String()
	actionCmdAddWebhook = actionCmdAdd.Flag("webhook", "The webhook to add the action to").HintAction(hintSubscriptionsNoNa).String()
	actionCmdAddCreate  = actionCmdAdd.Flag("create", "Create the file for the action").Default("false").Bool()
	//Action set
	actionCmdUpdate = actionCmd.Command("update", "Sets/Changes an action")
	//Action set webhook
	actionCmdUpdateWh        = actionCmdUpdate.Command("webhook", "Sets/Changes the webhook for an action")
	actionCmdUpdateWhAction  = actionCmdUpdateWh.Arg("action", "The action to change the webhook for").HintAction(hintListActions).Required().String()
	actionCmdUpdateWhWebhook = actionCmdUpdateWh.Arg("webhook", "The new webhook").HintAction(hintSubscriptions).Required().String()
	//Action create file
	actionCmdCreateFile       = actionCmd.Command("createFile", "Creates the file for an action")
	actionCmdCreateFileAction = actionCmdCreateFile.Arg("action", "The action to create the config for").HintAction(hintListActions).Required().String()
	//Action set file
	actionCmdUpdateAction      = actionCmdUpdate.Command("action", "Sets/Changes the webhook for an action")
	actionCmdUpdateFileAction  = actionCmdUpdateAction.Arg("action", "The action to change").HintAction(hintListActions).Required().String()
	actionCmdUpdateFileType    = actionCmdUpdateAction.Flag("new-mode", "The new kind of action. Leave empty to keep current value").HintAction(hintAvailableActions).String()
	actionCmdUpdateFileNewFile = actionCmdUpdateAction.Flag("new-file", "The new action-file").HintAction(hintListCurDir).String()
	//Action delete
	actionCmdDelete     = actionCmd.Command("delete", "Deletes an action from a webhook").Alias("rm")
	actionCmdDeleteName = actionCmdDelete.Arg("name", "The name of the action").HintAction(hintListActions).Required().Strings()

	//Source commands
	//Sources
	sourceCmd = app.Command("source", "Source command")
	//Sources
	sourcesCmd    = app.Command("sources", "List your sources")
	sourcesCmdAID = sourcesCmd.Arg("source", "View one specific source").String()
	sourcesCmdID  = sourcesCmd.Flag("source", "View one specific source").String()
	//Create source
	sourceCmdCreate            = sourceCmd.Command("create", "Create a new source").Alias("add")
	sourceCmdCreateName        = sourceCmdCreate.Arg("name", "The name of the source").Required().String()
	sourceCmdCreateMode        = sourceCmdCreate.Flag("mode", "The mode of the source").HintAction(hintAvailableActionsForSource).Required().String()
	sourceCmdCreateDescription = sourceCmdCreate.Arg("description", "The description of the source").String()
	sourceCmdCreatePrivate     = sourceCmdCreate.Flag("private", "If the source should be private").Default("false").Bool()
	//Delete source
	sourceCmdDelete   = sourceCmd.Command("delete", "Delete a source").Alias("rm").Alias("del")
	sourceCmdDeleteID = sourceCmdDelete.Arg("sourceID", "The ID of the source to delete").String()

	loginCmd     = app.Command("login", "login")
	loginCmdUser = loginCmd.Flag("username", "Your username").String()

	registerCmd = app.Command("register", "Create an account")
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
		db       *dbhelper.DBhelper
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
	case subscriptionsCmd:
		{
			//whsub subscriptions
			ViewSubscriptions(db, config)
		}
	case subscribeDelete.FullCommand():
		{
			//whsub subscription unsubscribe
			Unsubscribe(config, db, *subscribeDeleteID)
		}
	case subscribeImport.FullCommand():
		{
			//whsub subscription import
			ImportSubscription(db, *subscribeImportID)
		}

	//Actions --------------------
	case actionCmdAdd.FullCommand():
		{
			//whsub action add
			AddAction(db, *actionCmdAddMode, *actionCmdAddName, *actionCmdAddWebhook, *actionCmdAddFile, *actionCmdAddCreate)
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
	case actionCmdUpdateWh.FullCommand():
		{
			//whsub	action update webhook
			ActionSetWebhook(db, *actionCmdUpdateWhWebhook, *actionCmdUpdateWhAction)
		}
	case actionCmdUpdateAction.FullCommand():
		{
			//whsub action update action
			ActionSetFile(db, *actionCmdUpdateFileAction, *actionCmdUpdateFileType, *actionCmdUpdateFileNewFile)
		}
	case actionCmdCreateFile.FullCommand():
		{
			//whsub action create
			ActionCreateFileFromName(db, *actionCmdCreateFileAction)
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
			CreateSource(config, *sourceCmdCreateName, *sourceCmdCreateDescription, *sourceCmdCreateMode, *sourceCmdCreatePrivate)
		}
	case sourceCmdDelete.FullCommand():
		{
			//whsub source delete
			DeleteSource(db, config, *sourceCmdDeleteID)
		}
	case sourcesCmd.FullCommand():
		{
			//whsub sources
			ListSources(db, config, *sourcesCmdID, *sourcesCmdAID)
		}

	//User
	case loginCmd.FullCommand():
		{
			//whsub login
			LoginCommand(config, *loginCmdUser)
		}
	case registerCmd.FullCommand():
		{
			//whsub register
			RegisterCommand(config)
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
