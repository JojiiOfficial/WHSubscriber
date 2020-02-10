package main

import (
	"database/sql"
	"strconv"

	godbhelper "github.com/JojiiOfficial/GoDBHelper"
)

//Tables
const (
	//TableSubscriptions table for subscriptions
	TableSubscriptions = "subscriptions"
	//TableActions table for actions
	TableActions = "actions"
)

func getInitSQL() godbhelper.QueryChain {
	return godbhelper.QueryChain{
		Name:  "initChain",
		Order: 0,
		Queries: godbhelper.CreateInitVersionSQL(
			godbhelper.InitSQL{
				Query:   "CREATE TABLE %s (`pkID` INTEGER PRIMARY KEY AUTOINCREMENT, `hookID` TEXT, `hookName` TEXT)",
				FParams: []string{TableSubscriptions},
			},
			godbhelper.InitSQL{
				Query:   "CREATE TABLE %s (`pkID` INTEGER PRIMARY KEY AUTOINCREMENT, `name` TEXT, `hookID` INTEGER, `mode` INTEGER, `file` TEXT)",
				FParams: []string{TableActions},
			},
		),
	}
}

//Global sql queries

// ---------------------------------- Inserts ------------------------------------

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

func (action *Action) insert(dab *godbhelper.DBhelper) error {
	var rs *sql.Result
	var err error
	var id int64

	if rs, err = dab.Insert(*action, TableActions); err != nil || rs == nil {
		return err
	}

	if id, err = (*rs).LastInsertId(); err != nil {
		return err
	}

	action.ID = id
	return nil
}

// ---------------------------------- Selects ------------------------------------

func getActions(dab *godbhelper.DBhelper) ([]Action, error) {
	var actions []Action
	err := dab.QueryRowsf(&actions, "SELECT * FROM %s ORDER BY pkID DESC", []string{TableActions})
	return actions, err
}

func getHooksHumanized(dab *godbhelper.DBhelper, limit int) ([]string, error) {
	var hooks []string
	err := dab.QueryRows(&hooks, "SELECT hookName || '-' || hookID FROM "+TableSubscriptions)
	return hooks, err
}

func getActionIDs(dab *godbhelper.DBhelper, limit int) ([]string, error) {
	var hooks []string
	err := dab.QueryRowsf(&hooks, "SELECT pkID FROM %s LIMIT %s", []string{TableActions, strconv.Itoa(limit)})
	return hooks, err
}

func hasAction(dab *godbhelper.DBhelper, actionID int64) (bool, error) {
	var c int
	err := dab.QueryRowf(&c, "SELECT COUNT(*) FROM %s WHERE pkID=?", []string{TableActions}, actionID)
	if err != nil {
		return false, err
	}
	return c == 1, nil
}

// ---------------------------------- Deletions ------------------------------------

func deleteActionByID(dab *godbhelper.DBhelper, id int64) error {
	_, err := dab.Execf("DELETE FROM %s WHERE pkID=?", []string{TableActions}, id)
	return err
}
