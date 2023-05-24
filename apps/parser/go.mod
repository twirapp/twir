module github.com/satont/tsuwari/apps/parser

go 1.20

require github.com/tidwall/gjson v1.14.4

require (
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.0 // indirect
)

replace github.com/satont/tsuwari/libs/integrations/spotify => ../../libs/integrations/spotify

replace github.com/satont/tsuwari/libs/config => ../../libs/config

replace github.com/satont/tsuwari/libs/gomodels => ../../libs/gomodels

replace github.com/satont/tsuwari/libs/types => ../../libs/types

replace github.com/satont/tsuwari/libs/grpc => ../../libs/grpc

replace github.com/satont/tsuwari/libs/twitch => ../../libs/twitch

replace github.com/satont/tsuwari/libs/gopool => ../../libs/gopool
