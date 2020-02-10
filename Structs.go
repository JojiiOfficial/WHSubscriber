package main

//Webhook webhook struct
type Webhook struct {
	ID     int64  `db:"pkID" orm:"pk,ai"`
	HookID string `db:"hookID"`
}
