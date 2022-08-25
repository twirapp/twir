package random

import (
	"math/rand"
	"strconv"
	"strings"
	"time"
	"tsuwari/parser/internal/types"
)

const Name = "random"

func Handler(data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
	rand.Seed(time.Now().UnixNano())

	params := [2]int{0,50}
	if data.Params != nil {
		parsed := strings.Split(*data.Params, "-")

		first, err := strconv.Atoi(parsed[0])
		if err == nil {
			params[0] = first
		}
		second, err := strconv.Atoi(parsed[1])
		if err == nil {
			params[1] = second
		}
	}

	random := params[0] + rand.Intn(params[1] - params[0] + 1)
	result := types.VariableHandlerResult{Result: strconv.Itoa(random)}

	return &result, nil
}