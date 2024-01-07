package binaries

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/twirapp/twir/cli/internal/shell"
)

var binaries = []string{
	"google.golang.org/protobuf/cmd/protoc-gen-go@v1.31.0",
	"google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0",
	"github.com/twitchtv/twirp/protoc-gen-twirp@v8.1.3",
	"github.com/caddyserver/caddy/v2/cmd/caddy@v2.7.6",
}

func InstallGoBinaries() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	for _, bin := range binaries {
		splittedBinaryName := strings.Split(bin, "/")
		binaryName := strings.Split(splittedBinaryName[len(splittedBinaryName)-1], "@")[0]

		if isBinaryInstalled(binaryName) {
			continue
		}

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
