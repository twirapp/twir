package variables

import (
	"regexp"
	"tsuwari/parser/internal/types"
	"tsuwari/parser/internal/variables/random"
	"tsuwari/parser/internal/variables/sender"
)

var (
	Variables = make(map[string]types.Variable)
	Regexp = regexp.MustCompile(`\$\(([^)|]+)(?:\|([^)]+))?\)`)
)

func SetVariables() {
	Variables[random.Name] = types.Variable{
		Name: random.Name,
		Handler: random.Handler,
	}
	Variables[sender.Name] = types.Variable{
		Name: sender.Name,
		Handler: sender.Handler,
	}
}

func ParseVariables(input string) string {
	result := Regexp.ReplaceAllStringFunc(input, func(s string) string {
		v := Regexp.FindStringSubmatchIndex(s)
		matchedVarName := s[v[2]:v[3]]

		var params *string

		if v[4] != -1 {
			p := s[v[4]:v[5]]
			params = &p
		}

		if val, ok := Variables[matchedVarName]; ok {
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

	return result
}