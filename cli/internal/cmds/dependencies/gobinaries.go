package dependencies

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/twirapp/twir/cli/internal/shell"
)

var binaries = []string{
	"google.golang.org/protobuf/cmd/protoc-gen-go@v1.31.0",
	"google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0",
	"github.com/twitchtv/twirp/protoc-gen-twirp@v8.1.3",
	"github.com/caddyserver/caddy/v2/cmd/caddy@latest",
}

func installGoBinaries() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	for _, bin := range binaries {
		cmd, err := shell.CreateCommand(
			shell.ExecCommandOpts{
				Command: "go install " + bin,
				Pwd:     wd,
			},
		)
		if err != nil {
			return err
		}

		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, "GOBIN="+filepath.Join(wd, ".bin"), "GOFLAGS=")

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("cannot install %s: %w", bin, err)
		}
	}

	return nil
}
