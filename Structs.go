package main

//Webhook webhook struct
type Webhook struct {
	ID     int64  `db:"pkID" orm:"pk,ai"`
	HookID string `db:"hookID"`
}

//Action webhook struct
type Action struct {
	ID     int64  `db:"pkID" orm:"pk,ai"`
	Name   string `db:"name"`
	HookID string `db:"hookID"`
	Mode   int16  `db:"mode"`
	File   string `db:"file"`
}
