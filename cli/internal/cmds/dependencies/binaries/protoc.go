package binaries

import (
	"archive/zip"
	"bytes"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const protocVersion = "v29.0"
const protocUrl = "https://github.com/protocolbuffers/protobuf/releases/download"

func InstallProtoc() error {
	var arch string
	var platform string

	switch runtime.GOARCH {
	case "amd64":
		arch = "x86_64"
	case "amd32":
		arch = "x86_32"
	case "386":
		arch = "x86_32"
	case "arm64":
		arch = "aarch_64"
	}
	switch runtime.GOOS {
	case "darwin":
		platform = "osx"
	case "linux":
		platform = "linux"
	case "windows":
		platform = "win64"
	}

	if isBinaryInstalled("protoc") {
		return nil
	}

	var url strings.Builder
	url.WriteString(protocUrl)
	url.WriteString("/")
	url.WriteString(protocVersion)
	url.WriteString("/")
	url.WriteString("protoc-")
	url.WriteString(strings.Replace(protocVersion, "v", "", 1))
	url.WriteString("-")
	url.WriteString(platform)

	if platform != "win64" {
		url.WriteString("-")
		url.WriteString(arch)
	}

	url.WriteString(".zip")

	// https://github.com/protocolbuffers/protobuf/releases/download/v25.1/protoc-25.1-win64-x86_64.zip

	resp, err := http.Get(url.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	bodyReader := bytes.NewReader(bodyBytes)
	zipReader, err := zip.NewReader(bodyReader, int64(len(bodyBytes)))
	if err != nil {
		return err
	}

	for _, file := range zipReader.File {
		fileInfo := file.FileInfo()
		if fileInfo.IsDir() {
			continue
		}

		if !strings.Contains(file.Name, "bin/protoc") {
			continue
		}

		protocFile, err := file.Open()
		if err != nil {
			return err
		}

		protocBytes, err := io.ReadAll(protocFile)
		if err != nil {
			return err
		}

		wd, err := os.Getwd()
		if err != nil {
			return err
		}

		splittedFileName := strings.Split(file.Name, "/")
		fileName := splittedFileName[len(splittedFileName)-1]

		err = os.WriteFile(filepath.Join(wd, ".bin", fileName), protocBytes, 0755)
		if err != nil {
			return err
		}
	}

	return nil
}
