package main

import (
	"fmt"
)

func printServerVersion() {
	fmt.Printf("Server running on: v%.1f\n", ServerVersion)
}

func runWHReceiverServer() {
	fmt.Println("run the server")
}
