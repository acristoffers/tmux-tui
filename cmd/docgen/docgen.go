package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/acristoffers/tmux-tui/cmd"
	"github.com/spf13/cobra/doc"
)

func main() {
	genManPages("build/share/man/man1")
	genShellCompletions("build/share")
}

func genManPages(path string) {
	path = mkpath(path)

	header := &doc.GenManHeader{
		Title:   "tmux-tui",
		Section: "1",
	}

	if err := doc.GenManTree(cmd.RootCmd, header, path); err != nil {
		fmt.Fprintf(os.Stderr, "An error occurred: %s\n", err)
		os.Exit(1)
	}
}

func genShellCompletions(folder string) {
	bash := mkpath(path.Join(folder, "bash", "completions"))
	cmd.RootCmd.GenBashCompletionFileV2(path.Join(bash, "tmux-tui.bash"), true)

	fish := mkpath(path.Join(folder, "fish", "completions"))
	cmd.RootCmd.GenFishCompletionFile(path.Join(fish, "tmux-tui.fish"), true)

	zsh := mkpath(path.Join(folder, "zsh", "completions"))
	cmd.RootCmd.GenZshCompletionFile(path.Join(zsh, "_tmux-tui"))
}

func mkpath(path string) string {
	path, err := filepath.Abs(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "An error occurred: %s\n", err)
		os.Exit(1)
	}

	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		fmt.Fprintf(os.Stderr, "An error occurred: %s\n", err)
		os.Exit(1)
	}

	return path
}
