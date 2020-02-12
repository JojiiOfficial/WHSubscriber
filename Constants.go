package main

//Files
const (
	//DefaultConfigFile the default configuration file
	DefaultConfigFile = "config.yml"
	DataDir           = ".whsub"
)

//Remote endpoints
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

//Local endpoints
const (
	//Webhooks
	LEPWebhooks = "/hooks/"
)

var (
	//DefaultDatabaseFile database file
	DefaultDatabaseFile = "data.db"
)
