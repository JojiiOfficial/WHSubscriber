package main

//Subscription webhook struct
type Subscription struct {
	ID             int64  `db:"pkID" orm:"pk,ai"`
	Name           string `db:"hookName"`
	SourceID       string `db:"sourceID"`
	SubscriptionID string `db:"subsID"`
	IsValid        bool   `db:"valid"`
	Mode           uint8  `db:"mode"`
}

//Action webhook struct
//Mode (0 = script, 1 = action)
type Action struct {
	ID             int64  `db:"pkID" orm:"pk,ai"`
	Name           string `db:"name"`
	SubscriptionID int64  `db:"subscriptionID"`
	Mode           uint8  `db:"mode"`
	File           string `db:"file"`
	HookName       string `db:"hookName" orm:"-"`
}

// ------------------ Request structs ------------------

// -------> Sources
type sourceAddRequest struct {
	Name        string `json:"name"`
	Description string `json:"descr"`
	Private     bool   `json:"private"`
	Mode        uint8  `json:"mode"`
}

type sourceRequest struct {
	SourceID string `json:"sid,omitempty"`
	Content  string `json:"content,omitempty"`
}

// -------> User
type credentialsRequest struct {
	Username string `json:"username"`
	Password string `json:"pass"`
}

// -------> Subscriptions
type subscriptionRequest struct {
	SourceID    string `json:"sid"`
	CallbackURL string `json:"cbUrl"`
}

type unsubscribeRequest struct {
	SubscriptionID string `json:"sid"`
}

type subscriptionUpdateCallbackRequest struct {
	SubscriptionID string `json:"subID"`
	CallbackURL    string `json:"cbUrl"`
}
