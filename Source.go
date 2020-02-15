package main

import (
	"fmt"
	"strconv"

	godbhelper "github.com/JojiiOfficial/GoDBHelper"
	"github.com/fatih/color"
)

//CreateSource creates a new source for webhooks
func CreateSource(config *ConfigStruct, name, description string, private bool) {
	if !checkLoggedIn(config) {
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

	var respData sourceAddResponse
	response, err := RestRequest2(EPSourceCreate, req, &respData, config)
	if err != nil {
		fmt.Println("Err:", err.Error())
		return
	}

	if response.Status == ResponseSuccess {
		fmt.Println(color.HiGreenString("Success!"), fmt.Sprintf("Source create successfully.\nID:\t%s\nSecret:\t%s", respData.SourceID, respData.Secret))
	} else {
		fmt.Println("Err:", response.Message)
	}
}

//DeleteSource deletes a source
func DeleteSource(db *godbhelper.DBhelper, config *ConfigStruct, sourceID string) {
	if !checkLoggedIn(config) {
		return
	}

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

	response, err := RestRequest2(EPSources, req, &res, config)

	if err != nil || response.Status != ResponseSuccess {
		return []sourceResponse{}, err
	}

	if response.Status == ResponseSuccess {
		return res.Sources, nil
	}

	return []sourceResponse{}, nil
}

//SourceList lists your sources
func SourceList(db *godbhelper.DBhelper, config *ConfigStruct, id string) {
	if !checkLoggedIn(config) {
		return
	}
	sources, err := getSources(db, config, id)
	if err != nil {
		fmt.Println("Err:", err.Error())
		return
	}

	if len(sources) == 1 {
		source := sources[0]
		fmt.Println(color.HiGreenString("ID:\t ") + source.SourceID)
		fmt.Println(color.HiGreenString("Name:\t ") + source.Name)
		if source.Description != "NULL" {
			fmt.Println(color.HiGreenString("Descr.:\t ") + source.Description)
		}
		if len(source.Secret) > 0 {
			fmt.Println(color.HiGreenString("Secret:\t ") + source.Secret)
		}
		fmt.Println(color.HiGreenString("Private: ") + strconv.FormatBool(source.IsPrivate))
	} else if len(sources) > 1 {
		headingColor := color.New(color.FgHiGreen, color.Underline, color.Bold)
		headingColor.Println("SourceID\t\t\t\tName\t\t\tCreation\t\tSecret")
		for _, source := range sources {
			name := source.Name
			if len(name) < 8 {
				name += "\t"
			}
			if len(name) < 12 {
				name += "\t"
			}
			fmt.Printf("%s\t%s\t%s\t%s\n", source.SourceID, name, source.CreationTime, source.Secret)
		}
	} else {
		fmt.Println("You don't have sources")
	}
}
