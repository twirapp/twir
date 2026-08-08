package main

import (
	"log"
)

func main() {
	application, err := initializeApplication()
	if err != nil {
		log.Fatal(err)
	}

	if err := application.Run(); err != nil {
		log.Fatal(err)
	}
}
