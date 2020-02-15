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
	http.HandleFunc("/ping", pingHandler)

	//Start the server
	if config.Server.UseTLS {
		//Start TLS server in background
		go (func() {
			log.Fatal(http.ListenAndServeTLS(config.Server.ListenAddress, config.Server.SSLCert, config.Server.SSLKey, nil))
		})()
		if debug {
			log.Printf("Started HTTPS server on address %s\n", config.Server.ListenAddress)
		}
	} else {
		//Start HTTP server in background
		go (func() {
			log.Fatal(http.ListenAndServe(config.Server.ListenAddress, nil))
		})()
		if debug {
			log.Printf("Started HTTP server on address %s\n", config.Server.ListenAddress)
		}
	}

	//keep program running
	for {
		time.Sleep(1 * time.Hour)
	}
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ll")
	ip := gaw.GetIPFromHTTPrequest(r)
	match, err := matchIPHost(ip, configs.Client.ServerURL)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if !match {
		w.WriteHeader(http.StatusNotImplemented)
		fmt.Fprintf(w, "not implemented")
		return
	}
	hookSourceID := r.Header.Get(HeaderSource)
	hookSubsID := r.Header.Get(HeaderSubsID)
	fmt.Println(hookSourceID, hookSubsID)
	if len(hookSourceID) > 0 && len(hookSubsID) > 0 {
		has, err := hasSubscribted(dbs, hookSubsID, hookSourceID)
		if err != nil {
			fmt.Println(err.Error())
		} else if has {
			err = validateSubscription(dbs, hookSubsID)
			if err != nil {
				fmt.Println(err.Error())
			}
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "OK")
			fmt.Println("Ping success")
			return
		}
	}
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, "not implemented")
}

//OnWebhookReceived
func webhookPage(w http.ResponseWriter, r *http.Request) {
	hookSource := r.Header.Get(HeaderSource)
	hookReceivedTime := r.Header.Get(HeaderReceived)

	if len(hookSource) > 0 && len(hookReceivedTime) > 0 {

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
		if err != nil {
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
		log.Printf("Found request without correct headers (%s,%s) from %s\n", hookSource, hookReceivedTime, gaw.GetIPFromHTTPrequest(r))
		sendNotSubscripted(w)
	}
}

func sendNotSubscripted(w http.ResponseWriter) {
	w.WriteHeader(http.StatusTeapot)
	fmt.Fprintf(w, "Don't send me those dumb fucking requests!")
}

func handleValidWebhook(c chan bool, subscription *Subscription, actions []Action, r *http.Request) {
	if len(actions) > 0 {
		//Read input
		b, err := ioutil.ReadAll(io.LimitReader(r.Body, 100000))
		if err != nil {
			log.Println(err.Error())
			return
		}
		r.Body.Close()
		c <- true

		//Create temp file
		tmpFile, err := ioutil.TempFile(os.TempDir(), "whsubaction-")
		if err != nil {
			log.Println("Cannot create temporary file:", err)
			return
		}
		//Write tempfile
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

			fmt.Println("temp file:", file)
			action.Run(file)
		}
	} else {
		c <- true
	}
}

//
// ----------------- API SERVER ---------------------
//

//StartAPIServer starts the intern api server
func StartAPIServer(config *ConfigStruct) {

}
