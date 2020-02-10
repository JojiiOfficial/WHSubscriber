package main

import (
	"fmt"
	"log"
)

func printSubscriberVersion() {
	fmt.Printf("WebHook subscriber version: v%.1f\n", ServerVersion)
}

func subscribe() {
	url := *subscribeWhURL
	whID := *subscribeWhID

	if len(url) == 0 {
		log.Fatalln("Callback url is empty!")
		return
	}

	fmt.Printf("Trying to subscribe to '%s' using callback URL '%s'\n", whID, url)
}
