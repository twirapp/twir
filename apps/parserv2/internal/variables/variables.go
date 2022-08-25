package variables

import (
	"tsuwari/parser/internal/types"
	"tsuwari/parser/internal/variables/random"
)

var Variables = make(map[string]types.Variable)
 
func SetVariables() {
	Variables[random.Name] = types.Variable{
		Name: random.Name,
		Handler: random.Handler,
	}
}