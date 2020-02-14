package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"

	"github.com/fatih/color"
)

//Remote endpoints
const (
	//Subscriptions
	EPSubscription         = "/sub"
	EPSubscriptionAdd      = EPSubscription + "/add"
	EPSubscriptionActivate = EPSubscription + "/activate"
	EPSubscriptionRemove   = EPSubscription + "/remove"

	//User
	EPUser       = "/user"
	EPUserCreate = EPUser + "/create"
	EPLogin      = "/login"

	//Source
	EPSource       = "/source"
	EPSourceCreate = EPSource + "/create"
	EPSourceInfa   = EPSource + "/info"
	EPSourceDelete = EPSource + "/delete"
)

//RestRequest requests
func RestRequest(file string, data interface{}, config *ConfigStruct) (string, error) {
	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: config.Client.IgnoreCert}}
	client := &http.Client{Transport: tr}

	//Build url
	u, err := url.Parse(config.Client.ServerURL)
	if err != nil {
		return "", err
	}
	u.Path = path.Join(u.Path, file)

	//Encode data
	bda, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	//Make request
	resp, err := client.Post(u.String(), "application/json", bytes.NewBuffer(bda))
	if err != nil {
		return "", err
	}

	//Read respons
	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(d), nil
}

func request(file string, inpData interface{}, outputdata interface{}, config *ConfigStruct) error {
	resp, err := RestRequest(file, inpData, config)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(resp), &outputdata)
}

//Returns true if success
func checkResponse(str string, arg ...string) bool {
	if str == ResponseErrorStr {
		if len(arg) > 0 {
			fmt.Print(arg[0])
		} else {
			fmt.Println(color.RedString("Error"), "doing request!")
		}
		return false
	}

	return str == ResponseSuccessStr
}
