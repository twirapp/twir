module github.com/satont/twir/libs/types

go 1.21.5

require (
	github.com/satont/twir/libs/gomodels v0.0.0-00010101000000-000000000000
	github.com/tkrajina/typescriptify-golang-structs v0.1.11
	github.com/twirapp/twir/libs/bus-core v0.0.0-00010101000000-000000000000
)

require (
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/guregu/null v4.0.0+incompatible // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/satori/go.uuid v1.2.0 // indirect
	github.com/stretchr/testify v1.8.4 // indirect
	github.com/tkrajina/go-reflector v0.5.6 // indirect
	gorm.io/gorm v1.25.7 // indirect
)

replace (
	github.com/satont/twir/libs/gomodels => ../gomodels
	github.com/twirapp/twir/libs/bus-core => ../bus-core
)
