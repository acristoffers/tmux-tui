package cmd

import (
	"fmt"
	"os"

	"github.com/acristoffers/tmux-tui/tmux_tui"
	"github.com/muesli/termenv"
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

		listThemes, err := cmd.Flags().GetBool("list-themes")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not parse options: %s\n", err)
			os.Exit(1)
		}

		if listThemes {
			fmt.Print("The available themes and the name you should use with the --theme option are:\n\n")
			for _, theme := range tmux_tui.AvailableThemes {
				fmt.Printf("%20s %s\n", tmux_tui.ThemeForName(theme).Name, theme)
			}
			return
		}

		themeHandle, err := cmd.Flags().GetString("theme")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not parse options: %s\n", err)
			os.Exit(1)
		}

		theme := tmux_tui.ThemeForName(themeHandle)
		if theme == nil {
			fmt.Fprintf(os.Stderr, "The selected theme does not exist. The available themes are:\n\n")
			for _, theme := range tmux_tui.AvailableThemes {
				fmt.Fprintf(os.Stderr, "%20s %s\n", tmux_tui.ThemeForName(theme).Name, theme)
			}
			os.Exit(1)
		}

		originalBackgroundColor := termenv.DefaultOutput().BackgroundColor()
		defer termenv.DefaultOutput().SetBackgroundColor(originalBackgroundColor)

		termenvBackgroundColor := termenv.ColorProfile().Color(string(theme.Background))
		termenv.DefaultOutput().SetBackgroundColor(termenvBackgroundColor)

		p := tmux_tui.NewApplication(*theme)
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
	RootCmd.Flags().StringP("theme", "t", "dracula", "Selects a theme. Default: dracula")
	RootCmd.Flags().Bool("list-themes", false, "Lists available themes.")
}
