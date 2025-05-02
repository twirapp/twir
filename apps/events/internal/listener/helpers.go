package listener

import (
	"fmt"
	"strings"

	"github.com/twirapp/twir/libs/bus-core/events"
)

func predictionMapTopPredictors(predictors []events.PredictionTopPredictor) string {
	mapped := make([]string, 0, len(predictors))

	for _, p := range predictors {
		if p.PointsWin == nil {
			mapped = append(
				mapped,
				fmt.Sprintf("%s - %v", p.UserName, p.PointsUsed),
			)
		} else {
			mapped = append(
				mapped,
				fmt.Sprintf("%s - %v(+%v)", p.UserName, p.PointsUsed, *p.PointsWin),
			)
		}
	}

	return strings.Join(mapped, " · ")
}
