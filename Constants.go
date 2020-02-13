package main

//Files
const (
	//DefaultConfigFile the default configuration file
	DefaultConfigFile = "config.yml"
	DataDir           = ".whsub"
)

var (
	//DefaultDatabaseFile database file
	DefaultDatabaseFile = "data.db"
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
	EPSourceInfa   = EPSource + "/info"
	EPSourceDelete = EPSource + "/delete"
)

//Local endpoints
const (
//Webhooks
)

const (
	//HeaderSource the sourceID of the incomming hook
	HeaderSource = "W_S_Source"
	//HeaderReceived the unixtime when the hook was received
	HeaderReceived = "W_S_Source"
)
