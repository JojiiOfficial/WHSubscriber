package main

import (
	"fmt"
	"log"

	dbhelper "github.com/JojiiOfficial/GoDBHelper"
	"github.com/fatih/color"
)

func printSubscriberVersion() {
	fmt.Printf("WebHook subscriber version: v%.1f\n", ServerVersion)
}

//Subscribe (config, hookCallbackURL, WebhookID)
func Subscribe(db *dbhelper.DBhelper, config *ConfigStruct, callbackURL, sourceID string) {
	remoteURL := config.Client.ServerURL

	if len(callbackURL) == 0 && len(config.Client.CallbackURL) == 0 {
		log.Fatalln("Callback url is empty!")
		return
	}

	if len(callbackURL) == 0 {
		callbackURL = config.Client.CallbackURL
	}

	if len(sourceID) != 32 {
		fmt.Println(color.HiRedString("Error:"), "SourceID invalid")
		return
	}

	if *appDebug {
		fmt.Printf("Trying to subscribe to '%s' using callback URL '%s' and remote '%s'\n", sourceID, callbackURL, remoteURL)
	}

	subscription := Subscription{
		SourceID: sourceID,
	}

	//Request subscription
	token := "-"
	if isLoggedIn(config) {
		token = config.User.SessionToken
	}

	removeInvalidSubs(db, sourceID)

	subsRequest := subscriptionRequest{
		CallbackURL: callbackURL,
		SourceID:    sourceID,
		Token:       token,
	}

	var subscrResponse subscriptionResponse
	response, err := RestRequest(EPSubscriptionAdd, subsRequest, &subscrResponse, config)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if response.Status == ResponseSuccess {
		subscription.SubscriptionID = subscrResponse.SubscriptionID
		subscription.Name = subscrResponse.Name
		subscription.Mode = subscrResponse.Mode

		err = subscription.insert(db)
		if err != nil {
			log.Printf(err.Error())
			return
		}
		fmt.Printf("%s subscribed to '%s'\n", color.HiGreenString("Successfully"), sourceID)
	} else {
		PrintResponseError(response)
	}

}

//Unsubscribe unsubscribe a subscription
func Unsubscribe(db *dbhelper.DBhelper, config *ConfigStruct, id string) {
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
		PrintResponseError(response)
	}
}

//SubscriptionUpdateCallback updates the callback of a subscription
func SubscriptionUpdateCallback(db *dbhelper.DBhelper, config *ConfigStruct, subscriptionID, newCallback string) {
	wdid, err := getWhIDFromHumanInput(db, subscriptionID)
	if err != nil {
		fmt.Println("Err:", err.Error())
		return
	}

	subscription, _ := getSubscriptionFromID(db, wdid)
	if len(newCallback) == 0 {
		newCallback = config.Client.CallbackURL
	}

	token := "-"
	if len(config.User.SessionToken) > 0 {
		token = config.User.SessionToken
	}

	request := subscriptionUpdateCallbackRequest{
		CallbackURL:    newCallback,
		SubscriptionID: subscription.SubscriptionID,
		Token:          token,
	}

	resp, err := RestRequest(EPSubscriptionUpdCallback, request, nil, config)
	if err != nil {
		fmt.Println("Err:", err.Error())
		return
	}

	if resp.Status == ResponseSuccess {
		fmt.Printf("%s updating callback URL to '%s'\n", color.HiGreenString("Success"), newCallback)
	} else {
		PrintResponseError(resp)
	}
}

//ViewSubscriptions views subscriptions
func ViewSubscriptions(db *dbhelper.DBhelper, config *ConfigStruct) {
	subscriptions, err := getSubscriptions(db)
	if err != nil {
		log.Fatalln(err.Error())
		return
	}
	headingColor := color.New(color.FgHiGreen, color.Underline, color.Bold)
	headingColor.Println("ID\tSourceID\t\t\t\tName")
	for _, subscription := range subscriptions {
		hookID := subscription.SourceID
		//var colorFunc color
		cs := color.New(color.FgHiWhite)
		if !subscription.IsValid {
			cs = color.New(color.FgMagenta)
		}

		cs.Printf("%d\t%s\t%s\n", subscription.ID, hookID, subscription.Name)
	}
	if len(subscriptions) == 0 {
		fmt.Println("No subscription yet")
	}
	fmt.Println()
}
