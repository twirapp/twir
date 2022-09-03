package customvar

import (
	CTX "context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
	"tsuwari/parser/internal/types"
	"tsuwari/parser/internal/variablescache"

	v8 "rogchap.com/v8go"
)

const Name = "customvar"

var Iso = v8.NewIsolate()
var Global = createGlobal(Iso)

type CustomVar struct {
	Type      *string `json:"type"`
	EvalValue *string `json:"evalValue"`
	Response  *string `json:"response"`
}

func Handler(ctx *variablescache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
	result := &types.VariableHandlerResult{}

	if data.Params == nil {
		return result, nil
	}

	v := getVarByName(ctx, *data.Params)

	if v == nil {
		return result, nil
	}

	if v.Type == nil {
		return result, nil
	}

	if *v.Type == "SCRIPT" {
		r := executeJs(*v.EvalValue)
		if r != nil {
			result.Result = *r
		}
	} else {
		result.Result = *v.Response
	}

	return result, nil
}

func getVarByName(ctx *variablescache.VariablesCacheService, name string) *CustomVar {
	variable := &CustomVar{}
	r, err := ctx.Services.Redis.Get(CTX.TODO(), fmt.Sprintf("variables:%s:%s", ctx.Context.ChannelId, name)).Result()
	if err == nil {
		json.Unmarshal([]byte(r), variable)
	}

	return variable
}

func executeJs(script string) *string {
	c := v8.NewContext(Iso, Global)
	valChannel := make(chan string, 1)

	go func() {
		c.RunScript(fmt.Sprintf("const result = async () => { %s }", script), "script.js")
		val, err := c.RunScript("result()", "script.js")

		if err == nil {
			p, _ := val.AsPromise()

			for p.State() == v8.Pending {
				continue
			}

			if p.State() == v8.Rejected {
				valChannel <- "error when requesting api."
			} else {
				valChannel <- p.Result().String()
			}
		} else {
			valChannel <- "internal error when executing script."
		}
	}()

	select {
	case val := <-valChannel:
		r := val
		return &r
	case <-time.After(500 * time.Millisecond):
		vm := c.Isolate()
		vm.TerminateExecution()
		r := "request not finished in 500ms"
		return &r
	}
}

func createGlobal(iso *v8.Isolate) *v8.ObjectTemplate {
	global := v8.NewObjectTemplate(iso)
	fn := v8.NewFunctionTemplate(Iso, func(info *v8.FunctionCallbackInfo) *v8.Value {
		args := info.Args()
		url := args[0].String()

		resolver, _ := v8.NewPromiseResolver(info.Context())

		go func() {
			res, _ := http.Get(url)
			body, _ := ioutil.ReadAll(res.Body)
			contentType := res.Header.Get("content-type")

			if strings.HasPrefix(contentType, "application/json") {
				var data map[string]interface{}
				err := json.Unmarshal([]byte(string(body)), &data)

				if err != nil {
					v, _ := v8.NewValue(Iso, string("error when requesting api"))
					resolver.Reject(v)
				}

				jsonStr, err := json.Marshal(data)

				if err == nil {
					r, _ := v8.NewValue(Iso, string(jsonStr))
					resolver.Resolve(r)
				} else {
					r, _ := v8.NewValue(Iso, "")
					resolver.Reject(r)
				}
			} else {
				v, _ := v8.NewValue(Iso, string(body))
				resolver.Resolve(v)
			}
		}()
		return resolver.GetPromise().Value
	})

	global.Set("fetch", fn, v8.ReadOnly)
	return global
}
