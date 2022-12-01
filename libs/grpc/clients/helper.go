package clients

import "fmt"

func createClientAddr(env, service string, port int) string {
	ip := service
	if env != "production" {
		ip = "127.0.0.1"
	}

	return fmt.Sprintf("%s:%v", ip, port)
}
