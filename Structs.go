package main

//Subscription webhook struct
type Subscription struct {
	ID             int64  `db:"pkID" orm:"pk,ai"`
	Name           string `db:"hookName"`
	SourceID       string `db:"sourceID"`
	SubscriptionID string `db:"subsID"`
}

//Action webhook struct
//Mode (0 = script, 1 = action)
type Action struct {
	ID             int64  `db:"pkID" orm:"pk,ai"`
	Name           string `db:"name"`
	SubscriptionID int64  `db:"subscriptionID"`
	Mode           int8   `db:"mode"`
	File           string `db:"file"`
	HookName       string `db:"hookName" orm:"-"`
}

// ------------------ Request structs ------------------

//-----> Requests
type sourceAddRequest struct {
	Token       string `json:"token"`
	Name        string `json:"name"`
	Description string `json:"descr"`
	Private     bool   `json:"private"`
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"pass"`
}

type subscriptionRequest struct {
	Token       string `json:"token"`
	SourceID    string `json:"sid"`
	CallbackURL string `json:"cburl"`
}

type unsubscribeRequest struct {
	SubscriptionID string `json:"sid"`
}

type listSourcesRequest struct {
	Token    string `json:"token"`
	SourceID string `json:"sid,omitempty"`
}

//-----> Responses

type loginResponse struct {
	Status string `json:"status"`
	Token  string `json:"token"`
}

type sourceAddResponse struct {
	Status   string `json:"status"`
	Secret   string `json:"secret"`
	SourceID string `json:"id"`
}

type subscriptionResponse struct {
	Status         string `json:"status"`
	Message        string `json:"message,omitempty"`
	SubscriptionID string `json:"sid"`
	Name           string `json:"name"`
}

//Status a REST response status
type Status struct {
	StatusCode    string `json:"statusCode"`
	StatusMessage string `json:"statusMessage"`
}

type listSourcesResponse struct {
	Status  string           `json:"status"`
	Sources []sourceResponse `json:"sources,omitempty"`
}

type sourceResponse struct {
	Name         string `db:"name" json:"name"`
	SourceID     string `db:"sourceID" json:"sourceID"`
	Description  string `db:"description" json:"description"`
	Secret       string `db:"secret" json:"secret"`
	CreationTime string `db:"creationTime" json:"crTime"`
	IsPrivate    bool   `db:"private"`
}
