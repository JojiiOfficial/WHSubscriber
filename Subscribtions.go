package main

import (
	"fmt"
	"log"

	godbhelper "github.com/JojiiOfficial/GoDBHelper"
	"github.com/fatih/color"
)

func printSubscriberVersion() {
	fmt.Printf("WebHook subscriber version: v%.1f\n", ServerVersion)
}

//Subscribe (config, hookCallbackURL, WebhookID)
func Subscribe(db *godbhelper.DBhelper, config *ConfigStruct, callbackURL, webhookID string) {
	remoteURL := config.Client.ServerURL

	if len(callbackURL) == 0 && len(config.Client.CallbackURL) == 0 {
		log.Fatalln("Callback url is empty!")
		return
	}

	if len(callbackURL) == 0 {
		callbackURL = config.Client.CallbackURL
	}

	if *appDebug {
		fmt.Printf("Trying to subscribe to '%s' using callback URL '%s' and remote '%s'\n", webhookID, callbackURL, remoteURL)
	}

	wh := Subscription{
		SourceID: webhookID,
	}

	//Request subscription
	token := "-"
	if isLoggedIn(config) {
		token = config.User.SessionToken
	}
	subsRequest := subscriptionRequest{
		CallbackURL: callbackURL,
		SourceID:    webhookID,
		Token:       token,
	}
	var subscrResponse subscriptionResponse
	err := request(EPSubscriptionAdd, subsRequest, &subscrResponse, config)
	if err != nil {
		fmt.Println("Err:", err.Error())
		return
	}

	if !checkResponse(subscrResponse.Status, "Error subscribing: "+subscrResponse.Message) {
		return
	}

	wh.SubscriptionID = subscrResponse.SubscriptionID
	wh.Name = subscrResponse.Name

	err = wh.insert(db)
	if err != nil {
		log.Printf(err.Error())
		return
	}
	fmt.Printf("%s subscribed to '%s'\n", color.HiGreenString("Successfully"), webhookID)
}

//Unsubscribe unsubscribe a subscription
func Unsubscribe(db *godbhelper.DBhelper, id string) {
	//TODO Unsubscribe
	fmt.Println(id)
}

//ImportSubscription import a subscription
func ImportSubscription(db *godbhelper.DBhelper, id string) {
	//TODO Import subscripion
	fmt.Println(id)
}

//ViewSubscriptions views subscriptions
func ViewSubscriptions(db *godbhelper.DBhelper) {
	subscriptions, err := getSubscriptions(db)
	if err != nil {
		log.Fatalln(err.Error())
		return
	}
	headingColor := color.New(color.FgHiGreen, color.Underline, color.Bold)
	headingColor.Println("ID\tHookID\t\t\tName")
	for _, subscription := range subscriptions {
		hookID := subscription.SourceID
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
