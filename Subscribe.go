package main

import (
	"fmt"
	"log"
)

func printSubscriberVersion() {
	fmt.Printf("WebHook subscriber version: v%.1f\n", ServerVersion)
}

func subscribe() {
	callbackURL := *subscribeWhACallbackURL
	webhookID := *subscribeWhAID
	remoteURL := config.Client.ServerURL

	if len(callbackURL) == 0 && len(config.Server.CallbackURL) == 0 {
		log.Fatalln("Callback url is empty!")
		return
	}

	if len(callbackURL) == 0 {
		callbackURL = config.Server.CallbackURL
	}

	if *appDebug {
		fmt.Printf("Trying to subscribe to '%s' using callback URL '%s' and remote '%s'\n", webhookID, callbackURL, remoteURL)
	}

}
