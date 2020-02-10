package main

import (
	"fmt"
	"log"

	"github.com/fatih/color"
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

	wh := Subscription{
		HookID: webhookID,
	}
	err := wh.insert(db)
	if err != nil {
		log.Printf(err.Error())
		return
	}
	fmt.Println(wh.ID)
}

func printSubsciptionList() {
	subscriptions, err := getSubscriptions(db)
	if err != nil {
		log.Fatalln(err.Error())
		return
	}
	headingColor := color.New(color.FgHiGreen, color.Underline, color.Bold)
	headingColor.Println("ID\tHookID\t\t\tName")
	for _, subscription := range subscriptions {
		hookID := subscription.HookID
		if len(hookID) < 8 {
			hookID += "\t"
		}
		fmt.Printf("%d\t%s\t\t%s\n", subscription.ID, hookID, subscription.Name)
	}
	if len(subscriptions) == 0 {
		fmt.Println("No subscription yet")
	}
	fmt.Println()
}
