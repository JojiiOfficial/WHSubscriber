package main

import (
	"os"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	app   = kingpin.New("whsub", "A WebHook subscriber")
	debug = app.Flag("debug", "Enable debug mode.").Bool()

	runServer    = app.Command("server", "Commands for the WH subscriber server")
	startServer  = runServer.Command("start", "Start the server")
	serverVerion = runServer.Flag("version", "Show the version of the version of the server").Bool()

	subscribeWh = app.Command("subscribe", "Subscribe to a webhook")
	subWhID     = subscribeWh.Arg("whid", "Which webhook you want to subscribe").Required().String()
)

func init() {

}

func main() {
	app.HelpFlag.Short('h')
	if checkServerVersion() {
		return
	}

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case startServer.FullCommand():
		runWHReceiverServer()
	case subscribeWh.FullCommand():
		subscribe()
	}
}

func checkServerVersion() bool {
	args := os.Args
	if len(args) == 3 && args[1] == runServer.FullCommand() && args[2] == "--version" {
		printServerVersion()
		return true
	}
	return false
}
