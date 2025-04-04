package main

import (
	"github.com/twirapp/twir/apps/chat-translator/internal/app"
	"go.uber.org/fx"
)

func main() {
	fx.New(app.App).Run()
}
