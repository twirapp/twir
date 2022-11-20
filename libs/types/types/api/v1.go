package api

import "github.com/satont/tsuwari/libs/types/types/api/modules"

type Modules struct {
	Youtube modules.YouTube
}

type V1 struct {
	MODULES Modules
}
