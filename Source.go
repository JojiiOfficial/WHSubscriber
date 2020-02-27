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
	response, err := EPSourceCreate.DoRequest(req, &respData, true, config)
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
		PrintResponseError(response)
	}
}

func getURLFromSource(config *ConfigStruct, sourceID, secret string) string {
	sURL, _ := gaw.ParseURL(config.Client.ServerURL)
	sURL.JoinPath(fmt.Sprintf("/webhook/post/%s/%s/", sourceID, secret))
	u := (url.URL)(*sURL)
	return u.String()
}

//DeleteSource deletes a source
func DeleteSource(db *dbhelper.DBhelper, config *ConfigStruct, sourceID string) {

	req := sourceRequest{
		SourceID: sourceID,
		Token:    config.User.SessionToken,
		Content:  "-",
	}

	SourceUpdateRequest(db, config, EPSourceDelete, req, func(response *RestResponse) {
		if response.Status == ResponseSuccess {
			id, err := getSubscriptionID(db, sourceID)
			if err == nil {
				removeActionSource(db, id)
				deleteSubscriptionByID(db, sourceID)
			}
			fmt.Println("Source deleted", color.HiGreenString("successfully"))
		} else {
			PrintResponseError(response)
		}
	})
}

//UpdateSourceDescription updates the source description
func UpdateSourceDescription(db *dbhelper.DBhelper, config *ConfigStruct, sourceID, newText string) {
	if len(newText) == 0 {
		newText = "-"
	}

	req := sourceRequest{
		SourceID: sourceID,
		Token:    config.User.SessionToken,
		Content:  newText,
	}

	SourceUpdateRequest(db, config, EPSourceChangeDesc, req, func(resp *RestResponse) {
		if resp.Status == ResponseSuccess {
			str := "updating"
			if newText == "-" {
				str = "removing"
			}

			fmt.Printf("%s %s source description\n", color.HiGreenString("Success"), str)
		} else {
			fmt.Println("Error:", resp.Message)
		}
	})
}

//UpdateSourceName updates the name of a given source
func UpdateSourceName(db *dbhelper.DBhelper, config *ConfigStruct, sourceID, newName string) {
	req := sourceRequest{
		Content:  newName,
		Token:    config.User.SessionToken,
		SourceID: sourceID,
	}

	SourceUpdateRequest(db, config, EPSourceRename, req, func(response *RestResponse) {
		if response.Status == ResponseSuccess {
			fmt.Printf("%s updating source to '%s'\n", color.HiGreenString("Success"), newName)
		} else {
			PrintResponseError(response)
		}
	})
}

//ToggleSourceAccessMode updates the name of a given source
func ToggleSourceAccessMode(db *dbhelper.DBhelper, config *ConfigStruct, sourceID string) {
	req := sourceRequest{
		Content:  "-",
		Token:    config.User.SessionToken,
		SourceID: sourceID,
	}

	SourceUpdateRequest(db, config, EPSourceToggleAccess, req, func(response *RestResponse) {
		if response.Status == ResponseSuccess {
			fmt.Printf("%s updating source to '%s'\n", color.HiGreenString("Success"), response.Message)
		} else {
			PrintResponseError(response)
		}
	})
}

//SourceUpdateRequest updates a source
func SourceUpdateRequest(db *dbhelper.DBhelper, config *ConfigStruct, endpoint Endpoint, requestStruct sourceRequest, respFunc func(*RestResponse)) {
	if !checkLoggedIn(config) {
		return
	}

	if len(requestStruct.SourceID) != 32 {
		fmt.Println(color.HiRedString("Error:"), "SourceID invalid")
		return
	}

	if len(requestStruct.Content) > 149 {
		fmt.Println(color.HiRedString("Error:"), "Content too long!")
		return
	}

	response, err := endpoint.DoRequest(requestStruct, nil, true, config)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	respFunc(response)
}

func getSources(db *dbhelper.DBhelper, config *ConfigStruct, args ...string) ([]sourceResponse, error) {
	sid := "-"
	if len(args) > 0 && len(args[0]) > 0 {
		sid = args[0]
	}

	req := sourceRequest{
		SourceID: sid,
		Token:    config.User.SessionToken,
		Content:  "-",
	}
	var res listSourcesResponse

	response, err := EPSources.DoRequest(req, &res, true, config)

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
		fmt.Println(color.HiGreenString("Private: ") + strconv.FormatBool(source.IsPrivate))
		if len(source.Secret) > 0 {
			fmt.Println(color.HiGreenString("Secret:\t ") + source.Secret)
			fmt.Println(color.HiGreenString("URL:\t ") + getURLFromSource(config, source.SourceID, source.Secret))
		}
	} else if len(sources) > 1 {
		headingColor := color.New(color.FgHiGreen, color.Underline, color.Bold)
		headingColor.Println("SourceID\t\t\t\tName\t\t\tCreated\t\t\tSecret")
		for _, source := range sources {
			name := source.Name
			if len(name) < 8 {
				name += "\t"
			}
			if len(name) < 14 {
				name += "\t"
			}
			fmt.Printf("%s\t%s\t%s\t%s\n", source.SourceID, name, source.CreationTime, source.Secret)
		}
	} else {
		fmt.Println("You don't have sources")
	}
	fmt.Println()
}
