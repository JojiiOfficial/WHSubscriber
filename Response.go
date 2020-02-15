package main

import (
	"errors"
)

var (
	//ErrorInvalidHeaders error if response doesn't contain the right headers
	ErrorInvalidHeaders = errors.New("Invalid response headers")
	//ErrorErrorStatus error if response doesn't contain the right headers
	ErrorErrorStatus = errors.New("Error status")
)

//ResponseStatus the status of response
type ResponseStatus uint8

const (
	//ResponseError if there was an error
	ResponseError ResponseStatus = 0
	//ResponseSuccess if the response is successful
	ResponseSuccess ResponseStatus = 1
)

const (
	//HeaderStatus headername for status in response
	HeaderStatus string = "rstatus"
	//HeaderStatusMessage headername for status in response
	HeaderStatusMessage string = "rmess"
)

//RestResponse the response of a rest call
type RestResponse struct {
	HTTPCode int
	Status   ResponseStatus
	Message  string
}

type loginResponse struct {
	Token string `json:"token"`
}

type sourceAddResponse struct {
	Status   string `json:"status"`
	Secret   string `json:"secret"`
	SourceID string `json:"id"`
}

type subscriptionResponse struct {
	Message        string `json:"message,omitempty"`
	SubscriptionID string `json:"sid"`
	Name           string `json:"name"`
}

type listSourcesResponse struct {
	Sources []sourceResponse `json:"sources,omitempty"`
}

type sourceResponse struct {
	Name         string `json:"name"`
	SourceID     string `json:"sourceID"`
	Description  string `json:"description"`
	Secret       string `json:"secret"`
	CreationTime string `json:"crTime"`
	IsPrivate    bool   `json:"isPrivate"`
}
