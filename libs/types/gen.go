package main

import (
	"fmt"
	"os"

	"github.com/tkrajina/typescriptify-golang-structs/typescriptify"
	"github.com/twirapp/twir/libs/types/types/api"
	apioverlays "github.com/twirapp/twir/libs/types/types/api/overlays"
	"github.com/twirapp/twir/libs/types/types/overlays"
)

func main() {
	files := []string{"api", "overlays"}
	for _, f := range files {
		_ = os.Remove(fmt.Sprintf("src/%s.ts", f))
	}

	apiConverter := typescriptify.New().
		Add(api.V1{}).
		AddEnum(apioverlays.AllPresets)
	apiConverter.CreateInterface = true
	apiConverter.CreateConstructor = false
	apiConverter.CreateFromMethod = false

	err := apiConverter.ConvertToFile("src/api.ts")
	if err != nil {
		panic(err.Error())
	}

	overlaysConverter := typescriptify.New().
		Add(overlays.DudesGrowRequest{}).
		Add(overlays.DudesLeaveRequest{}).
		Add(overlays.DudesUserSettings{}).
		AddEnum(overlays.AllDudesSpriteEnumValues)
	overlaysConverter.CreateInterface = true
	overlaysConverter.CreateConstructor = false
	overlaysConverter.CreateFromMethod = false
	err = overlaysConverter.ConvertToFile("src/overlays.ts")
	if err != nil {
		panic(err)
	}
}
