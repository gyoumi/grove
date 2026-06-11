package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// findTailwind locates the Tailwind standalone CLI: PATH first, then the
// grove cache, downloading it on first use so apps need no Node toolchain.
func findTailwind() (string, error) {
	if p, err := exec.LookPath("tailwindcss"); err == nil {
		return p, nil
	}

	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	bin := filepath.Join(cacheDir, "grove", tailwindBinName())
	if _, err := os.Stat(bin); err == nil {
		return bin, nil
	}
	if err := downloadTailwind(bin); err != nil {
		return "", err
	}
	return bin, nil
}

func tailwindBinName() string {
	if runtime.GOOS == "windows" {
		return "tailwindcss.exe"
	}
	return "tailwindcss"
}

func tailwindAsset() (string, error) {
	switch runtime.GOOS + "/" + runtime.GOARCH {
	case "linux/amd64":
		return "tailwindcss-linux-x64", nil
	case "linux/arm64":
		return "tailwindcss-linux-arm64", nil
	case "darwin/amd64":
		return "tailwindcss-macos-x64", nil
	case "darwin/arm64":
		return "tailwindcss-macos-arm64", nil
	case "windows/amd64":
		return "tailwindcss-windows-x64.exe", nil
	}
	return "", fmt.Errorf("no tailwind standalone build for %s/%s", runtime.GOOS, runtime.GOARCH)
}

func downloadTailwind(dest string) error {
	asset, err := tailwindAsset()
	if err != nil {
		return err
	}
	url := "https://github.com/tailwindlabs/tailwindcss/releases/latest/download/" + asset
	fmt.Fprintf(os.Stderr, "grove: downloading tailwind standalone CLI (one-time) from %s\n", url)

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed: %s", resp.Status)
	}

	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return err
	}
	tmp := dest + ".tmp"
	f, err := os.OpenFile(tmp, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o755)
	if err != nil {
		return err
	}
	if _, err := io.Copy(f, resp.Body); err != nil {
		f.Close()
		os.Remove(tmp)
		return err
	}
	f.Close()
	return os.Rename(tmp, dest)
}
