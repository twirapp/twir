package main

import (
	"fmt"
	"os"

	"github.com/satont/twir/libs/types/types/api"
	"github.com/satont/twir/libs/types/types/api/overlays"
	"github.com/tkrajina/typescriptify-golang-structs/typescriptify"
)

func main() {
	files := []string{"api"}
	for _, f := range files {
		_ = os.Remove(fmt.Sprintf("src/%s.ts", f))
	}

	apiConverter := typescriptify.New().
		Add(api.V1{}).
		AddEnum(overlays.AllPresets)
	apiConverter.CreateInterface = true
	apiConverter.CreateConstructor = false
	apiConverter.CreateFromMethod = false

	err := apiConverter.ConvertToFile("src/api.ts")
	if err != nil {
		panic(err.Error())
	}
}
