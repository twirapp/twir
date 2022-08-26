package main

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
	"tsuwari/parser/internal/config/cfg"
	"tsuwari/parser/internal/config/redis"
	"tsuwari/parser/internal/types"
	"tsuwari/parser/internal/variables"
)

func main() {
	cfg.LoadConfig(".")
	redis.Connect()
	// nats.Connect()

	redis.Rdb.Do(
		redis.RedisCtx,
		redis.Rdb.B().Set().Key("latestTest").Value(strconv.FormatInt(time.Now().Unix(), 10),
	).Build()).Error()

	variables.SetVariables()
	regexp := regexp.MustCompile(`\$\(([^)|]+)(?:\|([^)]+))?\)`)

	input := regexp.ReplaceAllStringFunc("$(sender) $(random|1-1000) qwe", func(s string) string {
		v := regexp.FindStringSubmatchIndex(s)
		matchedVarName := s[v[2]:v[3]]

		var params *string

		if v[4] != -1 {
			p := s[v[4]:v[5]]
			params = &p
		}

		if val, ok := variables.Variables[matchedVarName]; ok {
			res, err := val.Handler(types.VariableHandlerParams{
				Key: matchedVarName,
				Params: params,
			})


			if err != nil {
				return string(err.Error())
			} else {
				return res.Result
			}
		}

		return s
	})

	fmt.Println(input)

	defer redis.Rdb.Close()
}