package weather

import (
	"context"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/imroc/req/v3"
	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
)

type weatherResponse struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Weather []struct {
		Description string `json:"description"`
	} `json:"weather"`
	Main struct {
		Humidity int     `json:"humidity"`
		Temp     float32 `json:"temp"`
	} `json:"main"`
	Clouds struct {
		All int `json:"all"`
	} `json:"clouds"`
	Wind struct {
		Speed float32 `json:"speed"`
	} `json:"wind"`
	Sys struct {
		Country string `json:"country"`
	} `json:"sys"`
}

var Weather = &types.Variable{
	Name:         "weather",
	Description:  lo.ToPtr("Get weather from OpenWeatherMap. If command used with param, then param will be used as city name."),
	Example:      lo.ToPtr("weather|en|London"),
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
		// take default city from variable params
		if len(params) == 2 && params[1] != "" {
			query = params[1]
		}
		// take city from command params
		if parseCtx.Text != nil {
			query = *parseCtx.Text
		}

		if query == "" {
			result.Result = paramsError
			return result, nil
		}

		data := weatherResponse{}
		resp, err := req.
			R().
			SetQueryParams(
				map[string]string{
					"appid": apiKey,
					"lang":  lang,
					"units": "metric",
					"q":     query,
				},
			).
			SetSuccessResult(&data).
			SetHeader("Content-Type", "application/json").
			Get("https://api.openweathermap.org/data/2.5/weather")
		if err != nil {
			return nil, err
		}
		if resp.StatusCode != 200 && resp.StatusCode != 404 {
			result.Result = fmt.Sprintf("OpenWeatherMap API error: %d", resp.StatusCode)
			return result, nil
		}
		if resp.StatusCode == 404 {
			result.Result = "Location not found"
			return result, nil
		}

		weatherDescription := make([]string, len(data.Weather))
		for i, weather := range data.Weather {
			r, j := utf8.DecodeRuneInString(weather.Description)
			weatherDescription[i] = string(unicode.ToTitle(r)) + weather.Description[j:]
		}

		result.Result = fmt.Sprintf(
			"%s (%s): %s, 🌡️ %s°C, ☁️ %d%%, 💦 %d%%, 💨 %s m/sec",
			data.Name,
			data.Sys.Country,
			strings.Join(weatherDescription, ", "),
			fmt.Sprintf("%.1f", data.Main.Temp),
			data.Clouds.All,
			data.Main.Humidity,
			fmt.Sprintf("%.1f", data.Wind.Speed),
		)

		return result, nil
	},
}
