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
