package main

import "fmt"

const (
	//SubscriberVersion version of the WebHook subscriber
	SubscriberVersion = 0.1
)

func printSubscriberVersion() {
	fmt.Printf("WebHook subscriber version: v%.1f\n", ServerVersion)
}

func subscribe() {
	url := *subscribeWhURL
	whID := *subscribeWhID

	fmt.Printf("Trying to subscribe to '%s' using callback URL '%s'\n", whID, url)
}
