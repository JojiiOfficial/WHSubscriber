package main

import "fmt"

const (
	//ServerVersion the version of the server
	ServerVersion = 0.1
)

func printServerVersion() {
	fmt.Printf("Server running on: v%.1f\n", ServerVersion)
}

func runWHReceiverServer() {
	fmt.Println("run the server")
}
