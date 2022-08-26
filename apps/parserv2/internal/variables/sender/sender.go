package sender

import (
	"tsuwari/parser/internal/types"
)

const Name = "sender"

func Handler(data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
	result := types.VariableHandlerResult{Result: "Satont"}

	return &result, nil
}