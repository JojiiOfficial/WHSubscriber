package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

//Endpoint a remote url-path
type Endpoint string

//Remote endpoints
const (
	//User
	EPUser       Endpoint = "/user"
	EPUserCreate Endpoint = EPUser + "/create"
	EPLogin      Endpoint = EPUser + "/login"

	//Source
	EPSource             Endpoint = "/source"
	EPSources            Endpoint = EPSource + "s"
	EPSourceCreate       Endpoint = EPSource + "/create"
	EPSourceUpdate       Endpoint = EPSource + "/update"
	EPSourceDelete       Endpoint = EPSourceUpdate + "/delete"
	EPSourceChangeDesc   Endpoint = EPSourceUpdate + "/changedescr"
	EPSourceRename       Endpoint = EPSourceUpdate + "/rename"
	EPSourceToggleAccess Endpoint = EPSourceUpdate + "/toggleAccess"

	//Subscriptions
	EPSubscription            Endpoint = "/sub"
	EPSubscriptionAdd         Endpoint = EPSubscription + "/add"
	EPSubscriptionActivate    Endpoint = EPSubscription + "/activate"
	EPSubscriptionRemove      Endpoint = EPSubscription + "/remove"
	EPSubscriptionUpdCallback Endpoint = EPSubscription + "/updateCallback"
)

//DoRequest a better request method
func (endpoint Endpoint) DoRequest(payload interface{}, retVar interface{}, useSession bool, config *ConfigStruct) (*RestResponse, error) {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: config.Client.IgnoreCert,
			},
		},
	}

	//Build url
	u, err := url.Parse(config.Client.ServerURL)
	if err != nil {
		return nil, err
	}
	u.Path = path.Join(u.Path, string(endpoint))

	//Encode data
	bda, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	//Make request
	req, _ := http.NewRequest("POST", u.String(), bytes.NewBuffer(bda))
	req.Header.Set("Content-Type", "application/json")
	if useSession {
		req.Header.Set("Authorization", "Bearer "+config.User.SessionToken)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	//Read and validate headers
	statusStr := resp.Header.Get(HeaderStatus)
	statusMessage := resp.Header.Get(HeaderStatusMessage)

	if len(statusStr) == 0 {
		return nil, ErrorInvalidHeaders
	}
	statusInt, err := strconv.Atoi(statusStr)
	if err != nil || (statusInt > 1 || statusInt < 0) {
		return nil, ErrorInvalidHeaders
	}
	status := (ResponseStatus)(uint8(statusInt))

	response := &RestResponse{
		HTTPCode: resp.StatusCode,
		Message:  statusMessage,
		Status:   status,
	}

	//Only fill retVar if response was successful
	if status == ResponseSuccess && retVar != nil {
		//Read response
		d, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		//Parse response into retVar
		err = json.Unmarshal(d, &retVar)
		if err != nil {
			return nil, err
		}
	}

	return response, nil
}
