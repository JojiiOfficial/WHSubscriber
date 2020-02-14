package main

import (
	"fmt"

	godbhelper "github.com/JojiiOfficial/GoDBHelper"
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

func getSources(db *godbhelper.DBhelper, config *ConfigStruct, args ...string) ([]sourceResponse, error) {
	sid := ""
	if len(args) > 0 {
		sid = args[0]
	}

	req := listSourcesRequest{
		SourceID: sid,
		Token:    config.User.SessionToken,
	}
	var res listSourcesResponse
	err := request(EPSources, req, &res, config)
	if err != nil {
		return []sourceResponse{}, err
	}
	if !checkResponse(res.Status, "Err") {
		return []sourceResponse{}, nil
	}
	return res.Sources, nil
}

//SourceList lists your sources
func SourceList(db *godbhelper.DBhelper, config *ConfigStruct) {
	if !isLoggedIn(config) {
		fmt.Println("You need to be logged in to use this feature")
		return
	}
	sources, err := getSources(db, config)
	if err != nil {
		fmt.Println("Err:", err.Error())
		return
	}

	if len(sources) > 0 {
		headingColor := color.New(color.FgHiGreen, color.Underline, color.Bold)
		headingColor.Println("SourceID\t\t\t\tName\t\tCreation\t\tSecret")
		for _, source := range sources {
			fmt.Printf("%s\t%s\t%s\t%s\n", source.SourceID, source.Name, source.CreationTime, source.Secret)
		}
	} else {
		fmt.Println("You don't have sources")
	}
}
