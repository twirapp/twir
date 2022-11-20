package main

import (
	"os"

	"github.com/satont/tsuwari/libs/types/types"
	"github.com/tkrajina/typescriptify-golang-structs/typescriptify"
)

func main() {
	converter := typescriptify.New().Add(types.YoutubeSettings{})
	converter.CreateInterface = true
	converter.CreateConstructor = false
	converter.CreateFromMethod = false
	os.Remove("src/generated.ts")
	err := converter.ConvertToFile("src/generated.ts")
	if err != nil {
		panic(err.Error())
	}
}
