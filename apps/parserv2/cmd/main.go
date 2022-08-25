package main

import (
	"fmt"
	"regexp"
	"tsuwari/parser/internal/types"
	"tsuwari/parser/internal/variables"
)


func substr(input string, start int, length int) string {
	asRunes := []rune(input)
	
	if start >= len(asRunes) {
			return ""
	}
	
	if start+length > len(asRunes) {
			length = len(asRunes) - start
	}
	
	return string(asRunes[start : start+length])
}

func main() {
	variables.SetVariables()
	regexp := regexp.MustCompile(`\$\(([^)|]+)(?:\|([^)]+))?\)`)

	input := regexp.ReplaceAllStringFunc("$(random|1-1000) qweqweqwe $(random|1-100000000000000000000)", func(s string) string {
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
				return s
			} else {
				return res.Result
			}
		}

		return s
	})

	fmt.Println(string(input))
}