package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/satont/twir/libs/grpc/servers"
)

type App struct {
	Name  string
	Stack string
	Port  int
}

func main() {
	// order matters
	apps := []App{
		{Stack: "go", Name: "tokens", Port: servers.TOKENS_SERVER_PORT},
		{Stack: "go", Name: "timers", Port: servers.TIMERS_SERVER_PORT},
		{Stack: "go", Name: "events", Port: servers.EVENTS_SERVER_PORT},
		{Stack: "node", Name: "integrations", Port: servers.INTEGRATIONS_SERVER_PORT},
		{Stack: "go", Name: "emotes-cacher", Port: servers.EMOTES_CACHER_SERVER_PORT},
		{Stack: "go", Name: "parser", Port: servers.PARSER_SERVER_PORT},
		{Stack: "go", Name: "eventsub", Port: servers.EVENTSUB_SERVER_PORT},
		{Stack: "node", Name: "eval", Port: servers.EVAL_SERVER_PORT},
		{Stack: "go", Name: "bots", Port: servers.BOTS_SERVER_PORT},
		{Stack: "go", Name: "watched", Port: servers.WATCHED_SERVER_PORT},
		{Stack: "go", Name: "websockets", Port: servers.WEBSOCKET_SERVER_PORT},
		{Stack: "node", Name: "ytsr", Port: servers.YTSR_SERVER_PORT},
		{Stack: "go", Name: "api-twirp", Port: 3002},
		{Stack: "go", Name: "scheduler"},
		{Stack: "frontend", Name: "dashboard-new", Port: 3006},
		{Stack: "frontend", Name: "landing", Port: 3005},
		{Stack: "frontend", Name: "overlays", Port: 3008},
		{Stack: "frontend", Name: "public", Port: 3007},
	}

	var processes []*os.Process

	for _, app := range apps {
		fmt.Println("Starting " + app.Name)

		var command string
		if app.Stack == "go" {
			command = "reflex -s -r '\\.go$' -- go run ./cmd/main.go"
		} else if app.Stack == "node" {
			command = "reflex -s -r '\\.ts$' -- sh -c 'TS_NODE_TRANSPILE_ONLY=true node --loader ts-node/esm --enable-source-maps --trace-warnings --nolazy src/index.ts'"
		} else {
			command = "pnpm dev"
		}

		cmd := exec.Command(
			"sh",
			"-c",
			command,
		)

		if app.Stack == "frontend" {
			cmd.Dir = fmt.Sprintf("./frontend/%s", app.Name)
		} else {
			cmd.Dir = fmt.Sprintf("./apps/%s", app.Name)
		}

		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err := cmd.Start()
		if err != nil {
			panic(err)
		}

		processes = append(processes, cmd.Process)

		if app.Port == 0 {
			continue
		}

		if app.Stack == "frontend" {
			continue
		}

		for {
			conn, _ := net.DialTimeout("tcp", net.JoinHostPort("", fmt.Sprintf("%v", app.Port)), 5*time.Second)
			if conn != nil {
				conn.Close()
				break
			} else {
				time.Sleep(500 * time.Millisecond)
				fmt.Println("Waiting " + app.Name + " to be ready...")
			}
		}
	}

	mainSignals := make(chan os.Signal, 1)
	signal.Notify(mainSignals, syscall.SIGTERM, syscall.SIGINT)
	<-mainSignals
	for _, process := range processes {
		process.Signal(syscall.SIGTERM)
	}
}
