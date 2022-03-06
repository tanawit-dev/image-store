package main

import (
	"fmt"
	"log"
	"os"
)

const Port = 8080

func main() {
	server, err := InitializeApplication()
	if err != nil {
		fmt.Printf("failed to intialize app: %s \n", err)
		os.Exit(2)
	}

	_ = server.Run(fmt.Sprintf("localhost:%d", Port))

	log.Printf("server is running on port %d", Port)
}
