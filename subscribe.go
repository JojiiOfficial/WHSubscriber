package main

import "fmt"

func subscribe() {
	url := *subscribeWhURL
	whID := *subscribeWhID

	fmt.Printf("Trying to subscribe to '%s' using callback URL '%s'\n", whID, url)
}
