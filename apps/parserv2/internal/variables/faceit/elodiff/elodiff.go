package faceitelodiff

import (
	"strconv"
	"tsuwari/parser/internal/types"
	"tsuwari/parser/internal/variablescache"
)

const Name = "faceit.todayEloDiff"
const Description = "Faceit today elo earned"

func Handler(ctx *variablescache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
	result := &types.VariableHandlerResult{}

	matches := ctx.GetFaceitLatestMatches()
	diff := ctx.GetFaceitTodayEloDiff(matches)

	result.Result = strconv.Itoa(diff)

	return result, nil
}
