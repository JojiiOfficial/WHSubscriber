package main

import (
	"fmt"
	"log"
	"net/http"
)

func printServerVersion() {
	fmt.Printf("Server running on: v%.1f\n", ServerVersion)
}

// ------------------ Receiver SERVER ----------------

//StartReceiverServer starts the receiver server
func StartReceiverServer(config *ConfigStruct) {
	if config.Server.Enable {
		http.HandleFunc(LEPWebhooks, webhookPage)
		if config.Server.UseTLS {
			log.Fatal(http.ListenAndServeTLS(config.Server.ListenAddress, config.Server.SSLCert, config.Server.SSLKey, nil))
		} else {
			log.Fatal(http.ListenAndServe(config.Server.ListenAddress, nil))
		}
	} else {
		fmt.Printf("Error: You need to enable the server first: 'enabled: true' (in the config)")
	}
}

//OnWebhookReceived
func webhookPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

//
// ----------------- API SERVER ---------------------
//

//StartAPIServer starts the intern api server
func StartAPIServer(config *ConfigStruct) {

}
