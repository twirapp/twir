package tlds

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/imroc/req/v3"
	"go.uber.org/fx"
)

type TLDS struct {
	List []string
}

func New(lc fx.Lifecycle) (*TLDS, error) {
	tlds := &TLDS{}

	err := tlds.fetch()
	if err != nil {
		return nil, err
	}

	return tlds, nil
}

func (c *TLDS) fetch() error {
	res, err := req.R().
		SetRetryCount(10).
		SetRetryInterval(
			func(resp *req.Response, attempt int) time.Duration {
				return 2 * time.Second
			},
		).
		Get("https://data.iana.org/TLD/tlds-alpha-by-domain.txt")
	if err != nil {
		return fmt.Errorf("cannot get tlds %w", err)
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("cannot get tlds %v %v", res.StatusCode, res.ErrorResult())
	}

	resString := ""
	defer res.Body.Close()
	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("cannot read body %w", err)
	}

	resString = string(bytes)

	splittedTlds := strings.Split(resString, "\n")
	splittedTlds = splittedTlds[1 : len(splittedTlds)-1]

	for i, v := range splittedTlds {
		splittedTlds[i] = strings.ToLower(v)
	}

	c.List = splittedTlds

	return nil
}

// var TLDS = func() []string {
// 	req, err := req.R().
// 		SetRetryCount(10).
// 		SetRetryInterval(
// 			func(resp *req.Response, attempt int) time.Duration {
// 				return 2 * time.Second
// 			},
// 		).
// 		Get("https://data.iana.org/TLD/tlds-alpha-by-domain.txt")
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	if req.StatusCode != 200 {
// 		panic(errors.New("cannot get tlds"))
// 	}
//
// 	resString := ""
// 	defer req.Body.Close()
// 	bytes, err := ioutil.ReadAll(req.Body)
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	resString = string(bytes)
//
// 	splittedTlds := strings.Split(resString, "\n")
// 	splittedTlds = splittedTlds[1 : len(splittedTlds)-1]
//
// 	for i, v := range splittedTlds {
// 		splittedTlds[i] = strings.ToLower(v)
// 	}
//
// 	return splittedTlds
// }()
