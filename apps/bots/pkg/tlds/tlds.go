package tlds

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type TLDS struct {
	List []string
}

func New() (*TLDS, error) {
	tlds := &TLDS{}

	err := tlds.fetch()
	if err != nil {
		return nil, err
	}

	return tlds, nil
}

func (c *TLDS) fetch() error {
	var res *http.Response
	var err error

	for attempt := 0; attempt < 10; attempt++ {
		res, err = http.Get("https://data.iana.org/TLD/tlds-alpha-by-domain.txt")
		if err == nil && res.StatusCode == 200 {
			break
		}
		if res != nil {
			res.Body.Close()
		}
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		return fmt.Errorf("cannot get tlds %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return fmt.Errorf("cannot get tlds %v", res.StatusCode)
	}

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("cannot read body %w", err)
	}

	resString := string(bytes)

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
