package main

import (
	"fmt"
	"os"

	"github.com/satont/tsuwari/libs/types/types/api"
	"github.com/tkrajina/typescriptify-golang-structs/typescriptify"
)

func main() {
	files := []string{"api"}
	for _, f := range files {
		os.Remove(fmt.Sprintf("src/%s.d.ts", f))
	}

	apiConverter := typescriptify.New().Add(api.V1{})
	apiConverter.CreateInterface = true
	apiConverter.CreateConstructor = false
	apiConverter.CreateFromMethod = false

	err := apiConverter.ConvertToFile("src/api.d.ts")
	if err != nil {
		panic(err.Error())
	}
}
