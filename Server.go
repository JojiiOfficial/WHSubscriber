package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	goHelper "github.com/JojiiOfficial/GoAwesomeHelper"
)

func printServerVersion() {
	fmt.Printf("Server running on: v%.1f\n", ServerVersion)
}

// ------------------ Receiver SERVER ----------------

//StartReceiverServer starts the receiver server
func StartReceiverServer(config *ConfigStruct, debug bool) {
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

		fmt.Fprintf(w, "Welcome to the HomePage!")
		fmt.Println("Endpoint Hit: homePage")
	} else {
		log.Printf("Found request without correct headers from %s\n", goHelper.GetIPFromHTTPrequest(r))
	}
}

//
// ----------------- API SERVER ---------------------
//

//StartAPIServer starts the intern api server
func StartAPIServer(config *ConfigStruct) {

}
