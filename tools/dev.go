package main

import (
	"flag"
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

func migrate() {
	cmd := exec.Command(
		"sh",
		"-c",
		"go run main.go",
	)

	cmd.Dir = "./libs/migrations"
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}

func caddy(withCaddy *bool) {
	if withCaddy == nil || !*withCaddy {
		return
	}

	cmd := exec.Command(
		"sh",
		"-c",
		"caddy reverse-proxy --from twir.satont.localhost --to 127.0.0.1:3005",
	)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	go func() {
		err := cmd.Run()
		if err != nil {
			panic(err)
		}
	}()
}

func main() {
	withCaddy := flag.Bool("caddy", false, "start caddy server?")
	flag.Parse()

	migrate()
	caddy(withCaddy)

	// order matters
	apps := []App{
		{Stack: "go", Name: "tokens", Port: servers.TOKENS_SERVER_PORT},
		{Stack: "go", Name: "events", Port: servers.EVENTS_SERVER_PORT},
		{Stack: "go", Name: "emotes-cacher", Port: servers.EMOTES_CACHER_SERVER_PORT},
		{Stack: "go", Name: "parser", Port: servers.PARSER_SERVER_PORT},
		{Stack: "go", Name: "eventsub", Port: servers.EVENTSUB_SERVER_PORT},
		{Stack: "node", Name: "eval", Port: servers.EVAL_SERVER_PORT},
		{Stack: "go", Name: "bots", Port: servers.BOTS_SERVER_PORT},
		{Stack: "go", Name: "timers", Port: servers.TIMERS_SERVER_PORT},
		{Stack: "go", Name: "websockets", Port: servers.WEBSOCKET_SERVER_PORT},
		{Stack: "go", Name: "ytsr", Port: servers.YTSR_SERVER_PORT},
		{Stack: "node", Name: "integrations", Port: servers.INTEGRATIONS_SERVER_PORT},
		{Stack: "go", Name: "api", Port: 3002},
		{Stack: "go", Name: "scheduler", Port: servers.SCHEDULER_SERVER_PORT},
		// {Stack: "go", Name: "discord", Port: servers.DISCORD_SERVER_PORT},
		{Stack: "frontend", Name: "dashboard", Port: 3006},
		{Stack: "frontend", Name: "landing", Port: 3005},
		{Stack: "frontend", Name: "overlays", Port: 3008},
		{Stack: "frontend", Name: "public-page", Port: 3007},
	}

	var processes []*os.Process

	for _, app := range apps {
		fmt.Println("Starting " + app.Name)
		// nodemon --exec "go run ./cmd/main.go" --ext "go" --watch . --cwd ./apps/tokens
		var command string
		if app.Stack == "go" {
			command = fmt.Sprintf(
				`pnpm nodemon --exec "go run" --ext "go" --watch . --cwd ./apps/%s --signal SIGTERM cmd/main.go`,
				app.Name,
			)
		} else if app.Stack == "node" {
			command = fmt.Sprintf(
				`pnpm nodemon --exec "cross-env TS_NODE_TRANSPILE_ONLY=true node --loader ts-node/esm --enable-source-maps --trace-warnings --nolazy" --ext "ts" --watch . --cwd ./apps/%s --signal SIGTERM src/index.ts`,
				app.Name,
			)
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
			conn, _ := net.DialTimeout(
				"tcp",
				net.JoinHostPort("", fmt.Sprintf("%v", app.Port)),
				5*time.Second,
			)
			if conn != nil {
				conn.Close()
				break
			} else {
				time.Sleep(500 * time.Millisecond)
				fmt.Print(".")
			}
		}
	}

	mainSignals := make(chan os.Signal, 1)
	signal.Notify(mainSignals, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	fmt.Printf("\n%v signal recieved", <-mainSignals)
	for _, process := range processes {
		process.Signal(syscall.SIGTERM)
	}
}
