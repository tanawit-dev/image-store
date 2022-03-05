package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	server, err := InitializeApplication()
	if err != nil {
		fmt.Printf("failed to create event: %s \n", err)
		os.Exit(2)
	}

	err = server.Run("localhost:8080")
	if err != nil {
		log.Fatalln(err)
		return
	}

	log.Printf("server is running")
}
