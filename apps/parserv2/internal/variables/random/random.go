package random

import (
	"errors"
	"math/rand"
	"strconv"
	"strings"
	"time"
	types "tsuwari/parser/internal/types"
	variables_cache "tsuwari/parser/internal/variablescache"

	"github.com/samber/lo"
)

var Variable = types.Variable{
	Name:        "random",
	Description: lo.ToPtr("Random number from N to N"),
	Handler: func(ctx *variables_cache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
		rand.Seed(time.Now().UnixNano())

		params := [2]int{0, 50}
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

		if params[0] > params[1] {
			return nil, errors.New("first number cannot be larger then second")
		}

		if params[0] < 0 || params[1] < 0 {
			return nil, errors.New("numbers cannot be lower then 0")
		}

		random := params[0] + rand.Intn(params[1]-params[0]+1)
		result := types.VariableHandlerResult{Result: strconv.Itoa(random)}

		return &result, nil
	},
}
