package main

import "log"

func main() {
	application, err := initializeApplication()
	if err != nil {
		log.Fatalf("initialize application: %v", err)
	}

	if err := application.Run(); err != nil {
		log.Fatalf("run application: %v", err)
	}
}
