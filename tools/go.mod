module github.com/satont/twir/tools

go 1.21

replace github.com/satont/twir/libs/grpc => ../libs/grpc

require github.com/satont/twir/libs/grpc v0.0.0-00010101000000-000000000000

require github.com/kvz/logstreamer v0.0.0-20221024075423-bf5cfbd32e39 // indirect
