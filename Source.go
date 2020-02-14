package main

import (
	"fmt"

	"github.com/fatih/color"
)

//CreateSource creates a new source for webhooks
func CreateSource(config *ConfigStruct, name, description string, private bool) {
	if !isLoggedIn(config) {
		fmt.Println("You need to be logged in to create sources!")
		return
	}
	if len(description) == 0 {
		description = "NULL"
	}
	req := sourceAddRequest{
		Description: description,
		Name:        name,
		Private:     private,
		Token:       config.User.SessionToken,
	}

	var createResponse sourceAddResponse
	err := request(EPSourceCreate, req, &createResponse, config)
	if err != nil {
		fmt.Println("Err:", err.Error())
		return
	}

	if !checkResponse(createResponse.Status, "Error creating source!") {
		return
	}
	fmt.Println(color.HiGreenString("Success!"), fmt.Sprintf("Source create successfully.\nID:\t%s\nSecret:\t%s", createResponse.SourceID, createResponse.Secret))
}

//SourceInfo shows informations for a source
func SourceInfo() {
	//TODO SourceInfo
}

//DeleteSource deletes a source
func DeleteSource() {
	//TODO DeleteSource
}
