package main

//Webhook webhook struct
type Webhook struct {
	ID     int64  `db:"pkID" orm:"pk,ai"`
	HookID string `db:"hookID"`
	Name   string `db:"hookName"`
}

//Action webhook struct
//Mode (0 = script, 1 = action)
type Action struct {
	ID       int64  `db:"pkID" orm:"pk,ai"`
	Name     string `db:"name"`
	HookID   int64  `db:"hookID"`
	Mode     int8   `db:"mode"`
	File     string `db:"file"`
	HookName string `db:"hookName" orm:"-"`
}
