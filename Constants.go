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

const (
	//HeaderSubsID the sourceID of the incoming hook
	HeaderSubsID = "W_S_SubsID"
	//HeaderSource the sourceID of the incoming hook
	HeaderSource = "W_S_Source"
	//HeaderReceived the unix time when the hook was received
	HeaderReceived = "W_S_Source"
)
