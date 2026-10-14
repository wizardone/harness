package main

import (
	"fmt"
	"harness/providers"
	"os"
	"path/filepath"
	"strings"
)

// renderHeader returns the model in bold, with the provider, token limit and directory faint underneath.
//   - header: the provider details to show.
//   - dir: the directory the harness was started in; the home directory is shortened to "~".
func renderHeader(header providers.Header, dir string) string {
	details := []string{header.Provider}
	if header.MaxTokens > 0 {
		details = append(details, fmt.Sprintf("max tokens %d", header.MaxTokens))
	}
	details = append(details, shortenHome(dir))
	return headerStyle.Render(header.Model) + "\n" + hintStyle.Render(strings.Join(details, " · "))
}

// shortenHome replaces the home directory at the start of dir with "~".
//   - dir: the directory to shorten.
func shortenHome(dir string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		return dir
	}
	if dir == home || strings.HasPrefix(dir, home+string(filepath.Separator)) {
		return "~" + strings.TrimPrefix(dir, home)
	}
	return dir
}
