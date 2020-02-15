package main

import (
	"fmt"
	"log"
	"time"

	dbhelper "github.com/JojiiOfficial/GoDBHelper"
	"github.com/fatih/color"
)

func printSubscriberVersion() {
	fmt.Printf("WebHook subscriber version: v%.1f\n", ServerVersion)
}

//Subscribe (config, hookCallbackURL, WebhookID)
func Subscribe(db *dbhelper.DBhelper, config *ConfigStruct, callbackURL, webhookID string) {
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

	subSourceIDtemp = webhookID

	var subscrResponse subscriptionResponse
	response, err := RestRequest(EPSubscriptionAdd, subsRequest, &subscrResponse, config)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if response.Status == ResponseSuccess {
		wh.SubscriptionID = subscrResponse.SubscriptionID
		wh.Name = subscrResponse.Name

		err = wh.insert(db)
		if err != nil {
			log.Printf(err.Error())
			return
		}
		fmt.Printf("%s subscribed to '%s'\n", color.HiGreenString("Successfully"), webhookID)
	} else {
		fmt.Println(color.HiRedString("Error:"), response.Message)
	}

	//Reset subSourceTemp
	go (func() {
		<-time.After(9 * time.Second)
		if len(subSourceIDtemp) > 0 {
			subSourceIDtemp = ""
		}
	})()
}

//Unsubscribe unsubscribe a subscription
func Unsubscribe(config *ConfigStruct, db *dbhelper.DBhelper, id string) {
	wdid, err := getWhIDFromHumanInput(db, id)
	if err != nil {
		fmt.Println("Err:", err.Error())
		return
	}
	subscription, _ := getSubscriptionFromID(db, wdid)

	//Request
	req := unsubscribeRequest{
		SubscriptionID: subscription.SubscriptionID,
	}
	response, err := RestRequest(EPSubscriptionRemove, req, nil, config)
	if err != nil {
		fmt.Println("Err:", err.Error())
		return
	}

	if response.Status == ResponseSuccess {
		err = deleteSubscription(db, wdid)
		if err != nil {
			fmt.Println("Err:", err.Error())
			return
		}
		fmt.Println(color.HiGreenString("Successfully"), "unsubscribed from", subscription.Name)
	} else {
		fmt.Println("Error:", response.Message)
	}
}

//ImportSubscription import a subscription
func ImportSubscription(db *dbhelper.DBhelper, id string) {
	fmt.Println(id)
}

//ViewSubscriptions views subscriptions
func ViewSubscriptions(db *dbhelper.DBhelper, config *ConfigStruct) {
	subscriptions, err := getSubscriptions(db)
	if err != nil {
		log.Fatalln(err.Error())
		return
	}
	headingColor := color.New(color.FgHiGreen, color.Underline, color.Bold)
	headingColor.Println("ID\tHookID\t\t\t\t\tName")
	for _, subscription := range subscriptions {
		hookID := subscription.SourceID
		fmt.Printf("%d\t%s\t%s\n", subscription.ID, hookID, subscription.Name)
	}
	if len(subscriptions) == 0 {
		fmt.Println("No subscription yet")
	}
	fmt.Println()
}
