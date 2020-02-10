package main

func inPortValid(port uint16) bool {
	return port > 0 && port < 65535
}
