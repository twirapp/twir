package hydrator

import (
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/goccy/go-json"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/valyala/fasttemplate"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Opts struct {
	fx.In

	Db *gorm.DB
}

func New(opts Opts) *Hydrador {
	return &Hydrador{
		db: opts.Db,
	}
}

type Hydrador struct {
	db *gorm.DB
}

var variablesRegular = regexp.MustCompile(
	`(?m)\$\((?P<all>(?P<main>[^.)|]+)(?:\.[^)|]+)?)(?:\|(?P<params>[^)]+))?\)`,
)

func (c *Hydrador) HydrateStringWithData(channelId string, str string, data any) (string, error) {
	template := fasttemplate.New(str, "{", "}")

	bytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	m := make(map[string]any)

	if err = json.Unmarshal(bytes, &m); err != nil {
		return "", err
	}

	s := template.ExecuteFuncString(
		func(w io.Writer, tag string) (int, error) {
			splittedTag := strings.Split(tag, ".")
			if len(splittedTag) > 1 {
				val, ok := m[splittedTag[0]].(map[string]any)
				if !ok {
					// key not found in map
					// return 0, fmt.Errorf("key '%s' is not a map[string]interface{}", splittedTag[0])
					return w.Write([]byte(""))
				}

				v, ok := val[splittedTag[1]]
				if !ok {
					// key not found in map
					// return 0, fmt.Errorf("key '%s' is not found in map", splittedTag[1])
					return w.Write([]byte(""))
				}

				return w.Write([]byte(fmt.Sprint(v)))
			} else {
				val, ok := m[tag]
				if !ok {
					// not a found
					// return 0, fmt.Errorf("key '%s' is not found", tag)
					return w.Write([]byte(""))
				}

				return w.Write([]byte(fmt.Sprint(val)))
			}
		},
	)

	for _, match := range variablesRegular.FindAllString(s, len(s)) {
		variable := variablesRegular.FindStringSubmatch(match)
		if len(variable) < 4 {
			continue
		}
		t := variable[len(variable)-2]
		variableName := variable[len(variable)-1]

		if t != "customvar" {
			continue
		}

		dbVariable := &model.ChannelsCustomvars{}
		err := c.db.
			Where(`"channelId" = ? AND "name" = ?`, channelId, variableName).
			Find(dbVariable).Error
		if err != nil {
			continue
		}

		if dbVariable.Type == model.CustomVarScript {
			continue
		}

		s = strings.ReplaceAll(s, match, dbVariable.Response)
	}

	return s, nil
}
