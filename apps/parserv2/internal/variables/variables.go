package variables

import (
	"tsuwari/parser/internal/types"
	"tsuwari/parser/internal/variables/random"
	"tsuwari/parser/internal/variables/sender"
)

var Variables = make(map[string]types.Variable)
 
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