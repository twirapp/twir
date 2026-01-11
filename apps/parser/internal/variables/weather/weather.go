package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"unicode"
	"unicode/utf8"

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
	CommandsOnly: false,
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

		u, err := url.Parse("https://api.openweathermap.org/data/2.5/weather")
		if err != nil {
			return nil, err
		}
		q := u.Query()
		q.Set("appid", apiKey)
		q.Set("lang", lang)
		q.Set("units", "metric")
		q.Set("q", query)
		u.RawQuery = q.Encode()

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
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

		var data weatherResponse
		if err := json.Unmarshal(body, &data); err != nil {
			return nil, err
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
