package cmd

import (
	"fmt"
	"os"

	"github.com/acristoffers/tmux-tui/tmux_tui"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "tmux-tui [PATH]",
	Short: "Terminal User Interface for managing tmux'es windows and sessions",
	Long:  "Allows you to create, rename, move and delete tmux'es windows and sessions",
	Run: func(cmd *cobra.Command, args []string) {
		version, err := cmd.Flags().GetBool("version")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not parse options: %s\n", err)
			os.Exit(1)
		}

		if version {
			fmt.Printf("Version %s", tmux_tui.Version)
			return
		}

		p := tmux_tui.NewApplication()
		m, err := p.Run()
		if err != nil {
			os.Stderr.WriteString(fmt.Sprintf("There's been an error: %v\n", err))
			os.Exit(1)
		}

		switch m := m.(type) {
		case tmux_tui.AppModel:
			if len(m.Error) != 0 {
				os.Stderr.WriteString(m.Error + "\n")
				os.Exit(1)
			}
		}
	},
}

func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	RootCmd.Flags().BoolP("version", "v", false, "Prints the version.")
}
