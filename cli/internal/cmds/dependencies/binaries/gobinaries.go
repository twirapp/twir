package binaries

import (
	"github.com/twirapp/twir/cli/internal/cmds/dependencies/binaries/gobinary"
)

var binaries = []gobinary.GoBinary{
	{Url: "google.golang.org/protobuf/cmd/protoc-gen-go@v1.31.0"},
	{Url: "google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0"},
	{Url: "github.com/twitchtv/twirp/protoc-gen-twirp@v8.1.3"},
	{Url: "github.com/caddyserver/caddy/v2/cmd/caddy@v2.7.6"},
}

func InstallGoBinaries() error {
	for _, bin := range binaries {
		name, version := bin.GetNameAndVersionFromUrl()

		isInstalled, err := bin.IsInstalled()
		if err != nil {
			return err
		}

		if isInstalled {
			binaryVersion, err := bin.GetGolangBinaryVersion(name)
			if err != nil {
				return err
			}

			if binaryVersion == version {
				continue
			}
		}

		if err := bin.Install(); err != nil {
			return nil
		}
	}

	return nil
}
