package main

import (
	"database/sql"

	godbhelper "github.com/JojiiOfficial/GoDBHelper"
)

//Tables
const (
	//TableSubscriptions table for subscriptions
	TableSubscriptions = "subscriptions"
)

func getInitSQL() godbhelper.QueryChain {
	return godbhelper.QueryChain{
		Name:  "initChain",
		Order: 0,
		Queries: godbhelper.CreateInitVersionSQL(
			godbhelper.InitSQL{
				Query:   "CREATE TABLE %s (`pkID` INTEGER PRIMARY KEY AUTOINCREMENT, `hookID` TEXT)",
				FParams: []string{TableSubscriptions},
			},
		),
	}
}

//Global sql queries

func (webhook *Webhook) insert(dab *godbhelper.DBhelper) error {
	var rs *sql.Result
	var err error
	var id int64

	if rs, err = dab.Insert(*webhook, TableSubscriptions); err != nil || rs == nil {
		return err
	}

	if id, err = (*rs).LastInsertId(); err != nil {
		return err
	}

	webhook.ID = id
	return nil
}
