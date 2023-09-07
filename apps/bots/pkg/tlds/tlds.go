package tlds

var TLDS = []string{}

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
