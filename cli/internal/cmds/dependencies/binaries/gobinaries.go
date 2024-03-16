package binaries

import (
	"github.com/twirapp/twir/cli/internal/cmds/dependencies/binaries/gobinary"
	"golang.org/x/sync/errgroup"
)

var binaries = []gobinary.GoBinary{
	{Url: "google.golang.org/protobuf/cmd/protoc-gen-go@v1.31.0"},
	{Url: "google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0"},
	{Url: "github.com/twitchtv/twirp/protoc-gen-twirp@v8.1.3"},
	{Url: "github.com/caddyserver/caddy/v2/cmd/caddy@v2.7.6"},
	{Url: "github.com/bufbuild/buf/cmd/buf@v1.27.0"},
}

func InstallGoBinaries() error {
	var wg errgroup.Group

	for _, bin := range binaries {
		bin := bin

		wg.Go(
			func() error {
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
						return nil
					}
				}

				if err := bin.Install(); err != nil {
					return err
				}

				return nil
			},
		)
	}

	return wg.Wait()
}
