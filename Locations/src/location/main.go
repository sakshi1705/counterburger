package main

import (
	"os"
)

func main() {

	portNumber := os.Getenv("PORT")
	if len(portNumber) == 0 {
		portNumber = "8000"
	}

	server := NewServerConfiguration()
	server.Run(":" + portNumber)
}
