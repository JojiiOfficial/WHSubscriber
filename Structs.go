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
