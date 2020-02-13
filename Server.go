package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	gaw "github.com/JojiiOfficial/GoAw"
	godbhelper "github.com/JojiiOfficial/GoDBHelper"
)

func printServerVersion() {
	fmt.Printf("Server running on: v%.1f\n", ServerVersion)
}

// ------------------ Receiver SERVER ----------------

var dbs *godbhelper.DBhelper

//StartReceiverServer starts the receiver server
func StartReceiverServer(config *ConfigStruct, db *godbhelper.DBhelper, debug bool) {
	dbs = db
	//Always listen only on /
	http.HandleFunc("/", webhookPage)

	//Start the server
	if config.Server.UseTLS {
		go (func() {
			log.Fatal(http.ListenAndServeTLS(config.Server.ListenAddress, config.Server.SSLCert, config.Server.SSLKey, nil))
		})()
		if debug {
			log.Printf("Started HTTPS server on address %s\n", config.Server.ListenAddress)
		}
	} else {
		go (func() {
			log.Fatal(http.ListenAndServe(config.Server.ListenAddress, nil))
		})()
		if debug {
			log.Printf("Started HTTP server on address %s\n", config.Server.ListenAddress)
		}
	}
	for {
		time.Sleep(1 * time.Hour)
	}
}

//OnWebhookReceived
func webhookPage(w http.ResponseWriter, r *http.Request) {
	hookSource := r.Header.Get(HeaderSource)
	hookReceivedTime := r.Header.Get(HeaderReceived)
	if len(hookSource) > 0 && len(hookReceivedTime) > 0 {
		fmt.Println("header", hookSource)

		subscription, err := getSubscriptionFromID(dbs, hookSource)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				log.Printf("A probably valid source sent you a webhook you have not subscripted: %s\n", hookSource)
				sendNotSubscripted(w)
				return
			}
			log.Printf("Err %s\n", err.Error())
			return
		}

		actions, err := getActionsForSource(dbs, subscription.ID)
		if err != nil {
			log.Printf("Your subscription '%s' was triggered but no action was found\n", subscription.Name)
		} else {
			//handle this webhook in a thread
			go handleValidWebhook(subscription, actions, r)
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK")
	} else {
		log.Printf("Found request without correct headers from %s\n", gaw.GetIPFromHTTPrequest(r))
		sendNotSubscripted(w)
	}
}

func sendNotSubscripted(w http.ResponseWriter) {
	w.WriteHeader(http.StatusTeapot)
	fmt.Fprintf(w, "Don't send me those dumb fucking requests!")
}

func handleValidWebhook(subscription *Subscription, actions []Action, request *http.Request) {
	if len(actions) > 0 {
		for _, action := range actions {
			fmt.Println(action.Name, "-", action.File, "-", action.Mode)
		}
		fmt.Print("handling actions")
	}
}

//
// ----------------- API SERVER ---------------------
//

//StartAPIServer starts the intern api server
func StartAPIServer(config *ConfigStruct) {

}
