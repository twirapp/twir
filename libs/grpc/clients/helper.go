package clients

import "fmt"

func createClientAddr(env, service string, port int) string {
	ip := "127.0.0.1"
	if env != "production" {
		ip = service
	}

	return fmt.Sprintf("%s:%v", ip, port)
}
