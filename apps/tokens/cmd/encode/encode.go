package main

import (
	"fmt"
	"log"
	"os"

	"github.com/twirapp/twir/libs/crypto"
)

func main() {
	args := os.Args
	if len(args) < 4 {
		log.Fatal("Wrong number of arguments")
	}

	cipherKey := args[1]

	accessToken, err := crypto.Encrypt(args[2], cipherKey)
	if err != nil {
		log.Fatal(err)
	}
	refreshToken, err := crypto.Encrypt(args[3], cipherKey)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("accessToken: ", accessToken)
	fmt.Println("refreshToken: ", refreshToken)
}
