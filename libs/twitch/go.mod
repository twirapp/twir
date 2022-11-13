module github.com/satont/tsuwari/libs/twitch

go 1.19

require (
	github.com/satont/go-helix/v2 v2.7.18
	github.com/satont/tsuwari/libs/gomodels v0.0.0-20221112130747-e34f337ae946
	gorm.io/gorm v1.24.1
)

require (
	github.com/golang-jwt/jwt/v4 v4.4.2 // indirect
	github.com/guregu/null v4.0.0+incompatible // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/lib/pq v1.10.7 // indirect
	github.com/rogpeppe/go-internal v1.9.0 // indirect
	github.com/satori/go.uuid v1.2.0 // indirect
)

replace github.com/satont/tsuwari/libs/gomodels => ../../libs/gomodels
