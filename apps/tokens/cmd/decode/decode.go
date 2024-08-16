package main

import (
	"fmt"
	"log"
	"os"

	"github.com/satont/twir/libs/crypto"
)

func main() {
	args := os.Args
	if len(args) < 3 {
		log.Fatal("Wrong number of arguments")
	}

	cipherKey := args[1]

	encoded, err := crypto.Decrypt(args[2], cipherKey)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("encoded: ", encoded)
}
