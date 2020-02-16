package main

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"

	gaw "github.com/JojiiOfficial/GoAw"
	dbhelper "github.com/JojiiOfficial/GoDBHelper"
	"github.com/fatih/color"
)

//CreateSource creates a new source for webhooks
func CreateSource(config *ConfigStruct, name, description, mode string, private bool) {
	if !checkLoggedIn(config) {
		return
	}

	if mode == "custom" {
		mode = "script"
	}

	m, has := Modes[mode]
	if !has {
		fmt.Printf("Mode not known '%s'\n", mode)
		return
	}

	if len(description) == 0 {
		description = "NULL"
	}
	req := sourceAddRequest{
		Description: description,
		Name:        name,
		Private:     private,
		Mode:        m,
		Token:       config.User.SessionToken,
	}

	var respData sourceAddResponse
	response, err := RestRequest(EPSourceCreate, req, &respData, config)
	if err != nil {
		fmt.Println("Err:", err.Error())
		return
	}

	if response.Status == ResponseSuccess {
		fmt.Println(color.HiGreenString("Success!"), fmt.Sprintf("Source create successfully.\nID:\t%s\nSecret:\t%s", respData.SourceID, respData.Secret))
		if err == nil {
			fmt.Printf("\nWebhook URL: %s\n", getURLFromSource(config, respData.SourceID, respData.Secret))
		}
	} else {
		fmt.Println("Err:", response.Message)
	}
}

func getURLFromSource(config *ConfigStruct, sourceID, secret string) string {
	surl, _ := gaw.ParseURL(config.Client.ServerURL)
	surl.JoinPath(fmt.Sprintf("/webhook/post/%s/%s/", sourceID, secret))
	u := (url.URL)(*surl)
	return u.String()
}

//DeleteSource deletes a source
func DeleteSource(db *dbhelper.DBhelper, config *ConfigStruct, sourceID string) {
	if !checkLoggedIn(config) {
		return
	}

	req := sourceRequest{
		SourceID: sourceID,
		Token:    config.User.SessionToken,
	}
	response, err := RestRequest(EPSourceDelete, req, nil, config)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if response.Status == ResponseSuccess {
		id, err := getSubscriptionID(db, sourceID)
		if err == nil {
			removeActionSource(db, id)
			deleteSubscriptionByID(db, sourceID)
		}
		fmt.Println("Source deleted", color.HiGreenString("successfully"))
	} else {
		fmt.Println("Error:", response.Message)
	}
}

func getSources(db *dbhelper.DBhelper, config *ConfigStruct, args ...string) ([]sourceResponse, error) {
	sid := ""
	if len(args) > 0 {
		sid = args[0]
	}

	req := sourceRequest{
		SourceID: sid,
		Token:    config.User.SessionToken,
	}
	var res listSourcesResponse

	response, err := RestRequest(EPSources, req, &res, config)

	if err != nil {
		return []sourceResponse{}, err
	}

	if response.Status == ResponseSuccess {
		return res.Sources, nil
	} else if response.Status == ResponseError {
		return []sourceResponse{}, errors.New(response.Message)
	}

	return []sourceResponse{}, nil
}

//ListSources lists your sources
func ListSources(db *dbhelper.DBhelper, config *ConfigStruct, idFlag, idArg string) {
	if !checkLoggedIn(config) {
		return
	}

	//Allow to use flag or arg
	id := idFlag
	if len(id) == 0 {
		id = idArg
	}

	//Get Sources from server
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
		fmt.Println(color.HiGreenString("URL:\t ") + getURLFromSource(config, source.SourceID, source.Secret))
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
