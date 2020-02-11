package main

//Files
const (
	//DefaultConfigFile the default configuration file
	DefaultConfigFile = "config.yml"
	DataDir           = ".whsub"
)

//Endpoints
const (
	//Subscriptions
	EPSubscription         = "/sub"
	EPSubscriptionAdd      = EPSubscription + "/add"
	EPSubscriptionActivate = EPSubscription + "/activate"
	EPSubscriptionRemove   = EPSubscription + "/remove"

	//User
	EPUser       = "/user"
	EPUserCreate = EPUser + "/create"

	//Source
	EPSource       = "/source"
	EPSourceCreate = EPSource + "/create"
	EPSourceDelete = EPSource + "/delete"
)

var (
	//DefaultDatabaseFile database file
	DefaultDatabaseFile = "data.db"
)
