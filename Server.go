package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	gaw "github.com/JojiiOfficial/GoAw"
	dbhelper "github.com/JojiiOfficial/GoDBHelper"
)

func printServerVersion() {
	fmt.Printf("Server running on: v%.1f\n", ServerVersion)
}

// ------------------ Receiver SERVER ----------------

var (
	dbs     *dbhelper.DBhelper
	configs *ConfigStruct
)

//StartReceiverServer starts the receiver server
func StartReceiverServer(config *ConfigStruct, db *dbhelper.DBhelper, debug bool) {
	dbs = db
	configs = config

	//Always listen only on /
	http.HandleFunc("/", webhookPage)

	//Start the server
	if config.Webserver.HTTPS.Enabled {
		//Start TLS server in background
		go (func() {
			log.Fatal(http.ListenAndServeTLS(config.Webserver.HTTPS.ListenAddress, config.Webserver.HTTPS.CertFile, config.Webserver.HTTPS.KeyFile, nil))
		})()
		if debug {
			log.Printf("Started HTTPS server on address %s\n", config.Webserver.HTTPS.ListenAddress)
		}
	}

	if config.Webserver.HTTP.Enabled {
		//Start HTTP server in background
		go (func() {
			log.Fatal(http.ListenAndServe(config.Webserver.HTTP.ListenAddress, nil))
		})()
		if debug {
			log.Printf("Started HTTP server on address %s\n", config.Webserver.HTTP.ListenAddress)
		}
	}

	//keep program running
	for {
		time.Sleep(1 * time.Hour)
	}
}

//OnWebhookReceived
func webhookPage(w http.ResponseWriter, r *http.Request) {
	hookSource := r.Header.Get(HeaderSource)
	hookSubscriptionID := r.Header.Get(HeaderSubsID)
	hookReceivedTime := r.Header.Get(HeaderReceived)

	if len(hookSource) > 0 && len(hookReceivedTime) > 0 && len(hookSubscriptionID) > 0 {
		has, err := hasSubscribted(dbs, hookSubscriptionID, hookSource)
		if err != nil || !has {
			if err != nil {
				log.Println("Error receiving hook:", err.Error())
			}
			return
		}

		subscription, err := getSubscriptionFromSourceID(dbs, hookSource)
		if err != nil {
			//Send not-subscribed message to server if source was not found in database
			if err.Error() == "sql: no rows in result set" {
				log.Printf("A probably valid source sent you a webhook you have not subscripted: %s\n", hookSource)
				sendNotSubscripted(w)
				return
			}
			//Log on different error
			log.Printf("Err %s\n", err.Error())
			return
		}

		actions, err := getActionsForSource(dbs, subscription.ID)
		if err != nil || len(actions) == 0 {
			//if no action was found
			log.Printf("Your subscription '%s' was triggered but no action was found\n", subscription.Name)
		} else {
			//handle this webhook in a thread
			c := make(chan bool, 1)
			go handleValidWebhook(c, subscription, actions, r)
			<-c
		}

		//send OK
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK")
	} else {
		//Request from sth else than the WeServer
		log.Printf("Found request without correct headers (%s,%s,%s) from %s\n", hookSource, hookReceivedTime, hookSubscriptionID, gaw.GetIPFromHTTPrequest(r))
		sendNotSubscripted(w)
	}
}

func sendNotSubscripted(w http.ResponseWriter) {
	w.WriteHeader(http.StatusTeapot)
	fmt.Fprintf(w, "Don't send me those dumb fucking requests!")
}

func handleValidWebhook(c chan bool, subscription *Subscription, actions []Action, r *http.Request) {
	//Read input
	b, err := ioutil.ReadAll(io.LimitReader(r.Body, 100000))
	if err != nil {
		log.Println(err.Error())
		return
	}

	parsedWebhook := WebhookData{
		Header:  r.Header,
		Payload: string(b),
	}

	r.Body.Close()
	c <- true

	//Create temp file
	tmpFile, err := ioutil.TempFile(os.TempDir(), "WhSubAction-")
	if err != nil {
		log.Println("Cannot create temporary file:", err)
		return
	}

	//Write temp-file
	_, err = tmpFile.Write(b)
	if err != nil {
		log.Printf("Error writing temp file '%s': %s\n", tmpFile.Name(), err.Error())
		return
	}
	file := tmpFile.Name()

	//Run actions
	for _, action := range actions {
		if *appDebug {
			fmt.Println(action.Name, "-", action.File, "-", action.Mode)
		}

		action.Run(file, subscription, &parsedWebhook)
	}

	//Delete tempFile
	os.Remove(tmpFile.Name())
}

//
// ----------------- API SERVER ---------------------
//

//StartAPIServer starts the intern api server
func StartAPIServer(config *ConfigStruct) {

}
