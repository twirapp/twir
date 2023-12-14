package weather

import (
	"context"
	"fmt"
	"strings"

	"github.com/imroc/req/v3"
	"github.com/samber/lo"
	"github.com/satont/twir/apps/parser/internal/types"
)

type WeatherResponse struct {
	ID      int
	Name    string
	Weather []struct {
		Description string
	}
	Main struct {
		Humidity int
		Temp     float32
	}
	Clouds struct {
		All int
	}
	Wind struct {
		Speed float32
	}
	Sys struct {
		Country string
	}
}

var Weather = &types.Variable{
	Name:         "weather",
	Description:  lo.ToPtr("Get weather from OpenWeatherMap"),
	Example:      lo.ToPtr("weather|en"),
	CommandsOnly: true,
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		result := &types.VariableHandlerResult{}

		apiKey := parseCtx.Services.Config.OpenWeatherMapApiKey
		if apiKey == "" {
			result.Result = "No API key provided"
			return result, nil
		}

		paramsError := "Parameters are not specified. Example: $(weather|en), $(weather|ru|Moscow)"
		if variableData.Params == nil {
			result.Result = paramsError
			return result, nil
		}

		params := strings.Split(*variableData.Params, "|")
		if params == nil {
			result.Result = paramsError
			return result, nil
		}

		for i, str := range params {
			params[i] = strings.TrimSpace(str)
			if len(params[i]) == 0 {
				result.Result = paramsError
				return result, nil
			}
		}

		lang := "en"
		if len(params) >= 1 && params[0] != "" {
			lang = params[0]
		}

		query := ""
		if parseCtx.Text != nil {
			query = *parseCtx.Text
		}

		if len(params) == 2 && params[1] != "" {
			query = params[1]
		}

		if query == "" {
			result.Result = paramsError
			return result, nil
		}

		data := WeatherResponse{}
		resp, err := req.
			R().
			SetQueryParams(map[string]string{
				"appid": apiKey,
				"lang":  lang,
				"units": "metric",
				"q":     query,
			}).
			SetSuccessResult(&data).
			SetHeader("Content-Type", "application/json").
			Get("https://api.openweathermap.org/data/2.5/weather")

		fmt.Println(resp, err)

		if err != nil {
			return nil, err
		}

		if resp.StatusCode == 404 {
			result.Result = "Location not found"
			return result, nil
		}

		weatherDescription := ""
		for _, weather := range data.Weather {
			weatherDescription += strings.ToUpper(weather.Description[0:1]) + weather.Description[1:]
		}

		result.Result = fmt.Sprintf(
			"%s (%s): %s, 🌡️ %s°C, ☁️ %d%%, 💦 %d%%, 💨 %s m/sec",
			data.Name,
			data.Sys.Country,
			weatherDescription,
			fmt.Sprintf("%.1f", data.Main.Temp),
			data.Clouds.All,
			data.Main.Humidity,
			fmt.Sprintf("%.1f", data.Wind.Speed),
		)

		return result, nil
	},
}
