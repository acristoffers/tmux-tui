package tmux_tui

import (
	"errors"

	"github.com/charmbracelet/lipgloss"
)

type Theme struct {
	Name                string
	Handle              string
	Background          lipgloss.Color
	Foreground          lipgloss.Color
	Accent              lipgloss.Color
	Secondary           lipgloss.Color
	SelectionBackground lipgloss.Color
}

var AvailableThemeHandles = []string{
	"ayu-dark",
	"cobalt2",
	"dracula",
	"dracula-pro",
	"github-dark",
	"gruvbox-dark",
	"gruvbox-light",
	"material-palenight",
	"monokai",
	"night-owl",
	"nord",
	"oceanic-next",
	"one-dark",
	"one-light",
	"papercolor-dark",
	"papercolor-light",
	"solarized-dark",
	"solarized-light",
	"synthwave",
	"tokyo-night",
	"tomorrow-night",
}

var AvailableThemes = []Theme{
	AyuDarkTheme,
	Cobalt2Theme,
	DraculaProTheme,
	DraculaTheme,
	GitHubDarkTheme,
	GruvboxDarkTheme,
	GruvboxLightTheme,
	MaterialPalenightTheme,
	MonokaiTheme,
	NightOwlTheme,
	NordTheme,
	OceanicNextTheme,
	OneDarkTheme,
	OneLightTheme,
	PaperColorDarkTheme,
	PaperColorLightTheme,
	SolarizedDarkTheme,
	SolarizedLightTheme,
	Synthwave84Theme,
	TokyoNightTheme,
	TomorrowNightTheme,
}

func ThemeForName(name string) (Theme, error) {
	switch name {
	case "ayu-dark":
		return AyuDarkTheme, nil
	case "cobalt2":
		return Cobalt2Theme, nil
	case "dracula":
		return DraculaTheme, nil
	case "dracula-pro":
		return DraculaProTheme, nil
	case "github-dark":
		return GitHubDarkTheme, nil
	case "gruvbox-dark":
		return GruvboxDarkTheme, nil
	case "gruvbox-light":
		return GruvboxLightTheme, nil
	case "material-palenight":
		return MaterialPalenightTheme, nil
	case "monokai":
		return MonokaiTheme, nil
	case "night-owl":
		return NightOwlTheme, nil
	case "nord":
		return NordTheme, nil
	case "oceanic-next":
		return OceanicNextTheme, nil
	case "one-dark":
		return OneDarkTheme, nil
	case "one-light":
		return OneLightTheme, nil
	case "papercolor-dark":
		return PaperColorDarkTheme, nil
	case "papercolor-light":
		return PaperColorLightTheme, nil
	case "solarized-dark":
		return SolarizedDarkTheme, nil
	case "solarized-light":
		return SolarizedLightTheme, nil
	case "synthwave":
		return Synthwave84Theme, nil
	case "tokyo-night":
		return TokyoNightTheme, nil
	case "tomorrow-night":
		return TomorrowNightTheme, nil
	}
	return Theme{}, errors.New("Theme not found")
}

var DraculaTheme = Theme{
	Name:                "Dracula",
	Handle:              "dracula",
	Background:          lipgloss.Color("#282A36"),
	Foreground:          lipgloss.Color("#E3E3DE"),
	Accent:              lipgloss.Color("#50FA7B"),
	Secondary:           lipgloss.Color("#FFB86C"),
	SelectionBackground: lipgloss.Color("#BD93F9"),
}

var MonokaiTheme = Theme{
	Name:                "Monokai",
	Handle:              "monokai",
	Background:          lipgloss.Color("#272822"),
	Foreground:          lipgloss.Color("#F8F8F2"),
	Accent:              lipgloss.Color("#A6E22E"),
	Secondary:           lipgloss.Color("#FD971F"),
	SelectionBackground: lipgloss.Color("#49483E"),
}

var SolarizedDarkTheme = Theme{
	Name:                "Solarized Dark",
	Handle:              "solarized-dark",
	Background:          lipgloss.Color("#002B36"),
	Foreground:          lipgloss.Color("#839496"),
	Accent:              lipgloss.Color("#268BD2"),
	Secondary:           lipgloss.Color("#2AA198"),
	SelectionBackground: lipgloss.Color("#073642"),
}

var SolarizedLightTheme = Theme{
	Name:                "Solarized Light",
	Handle:              "solarized-light",
	Background:          lipgloss.Color("#FDF6E3"),
	Foreground:          lipgloss.Color("#657B83"),
	Accent:              lipgloss.Color("#268BD2"),
	Secondary:           lipgloss.Color("#2AA198"),
	SelectionBackground: lipgloss.Color("#EEE8D5"),
}

var NordTheme = Theme{
	Name:                "Nord",
	Handle:              "nord",
	Background:          lipgloss.Color("#2E3440"),
	Foreground:          lipgloss.Color("#D8DEE9"),
	Accent:              lipgloss.Color("#88C0D0"),
	Secondary:           lipgloss.Color("#81A1C1"),
	SelectionBackground: lipgloss.Color("#4C566A"),
}

var GruvboxDarkTheme = Theme{
	Name:                "Gruvbox Dark",
	Handle:              "gruvbox-dark",
	Background:          lipgloss.Color("#282828"),
	Foreground:          lipgloss.Color("#EBDBB2"),
	Accent:              lipgloss.Color("#FABD2F"),
	Secondary:           lipgloss.Color("#FE8019"),
	SelectionBackground: lipgloss.Color("#3C3836"),
}

var GruvboxLightTheme = Theme{
	Name:                "Gruvbox Light",
	Handle:              "gruvbox-light",
	Background:          lipgloss.Color("#FBF1C7"),
	Foreground:          lipgloss.Color("#3C3836"),
	Accent:              lipgloss.Color("#FABD2F"),
	Secondary:           lipgloss.Color("#FE8019"),
	SelectionBackground: lipgloss.Color("#EBDBB2"),
}

var TomorrowNightTheme = Theme{
	Name:                "Tomorrow Night",
	Handle:              "tomorrow-night",
	Background:          lipgloss.Color("#1D1F21"),
	Foreground:          lipgloss.Color("#C5C8C6"),
	Accent:              lipgloss.Color("#81A2BE"),
	Secondary:           lipgloss.Color("#F0C674"),
	SelectionBackground: lipgloss.Color("#373B41"),
}

var OneDarkTheme = Theme{
	Name:                "One Dark",
	Handle:              "one-dark",
	Background:          lipgloss.Color("#282C34"),
	Foreground:          lipgloss.Color("#ABB2BF"),
	Accent:              lipgloss.Color("#61AFEF"),
	Secondary:           lipgloss.Color("#E06C75"),
	SelectionBackground: lipgloss.Color("#3E4451"),
}

var OneLightTheme = Theme{
	Name:                "One Light",
	Handle:              "one-light",
	Background:          lipgloss.Color("#FAFAFA"),
	Foreground:          lipgloss.Color("#383A42"),
	Accent:              lipgloss.Color("#61AFEF"),
	Secondary:           lipgloss.Color("#E06C75"),
	SelectionBackground: lipgloss.Color("#E5E5E6"),
}

var OceanicNextTheme = Theme{
	Name:                "Oceanic Next",
	Handle:              "oceanic-next",
	Background:          lipgloss.Color("#1B2B34"),
	Foreground:          lipgloss.Color("#D8DEE9"),
	Accent:              lipgloss.Color("#6699CC"),
	Secondary:           lipgloss.Color("#EC5f67"),
	SelectionBackground: lipgloss.Color("#343D46"),
}

var DraculaProTheme = Theme{
	Name:                "Dracula Pro",
	Handle:              "dracula-pro",
	Background:          lipgloss.Color("#1E1F29"),
	Foreground:          lipgloss.Color("#F8F8F2"),
	Accent:              lipgloss.Color("#BD93F9"),
	Secondary:           lipgloss.Color("#FF79C6"),
	SelectionBackground: lipgloss.Color("#44475A"),
}

var AyuDarkTheme = Theme{
	Name:                "Ayu Dark",
	Handle:              "ayu-dark",
	Background:          lipgloss.Color("#0F1419"),
	Foreground:          lipgloss.Color("#E6E1CF"),
	Accent:              lipgloss.Color("#39BAE6"),
	Secondary:           lipgloss.Color("#FF8F40"),
	SelectionBackground: lipgloss.Color("#273747"),
}

var NightOwlTheme = Theme{
	Name:                "Night Owl",
	Handle:              "night-owl",
	Background:          lipgloss.Color("#011627"),
	Foreground:          lipgloss.Color("#D6DEEB"),
	Accent:              lipgloss.Color("#82AAFF"),
	Secondary:           lipgloss.Color("#7E57C2"),
	SelectionBackground: lipgloss.Color("#5F7E97"),
}

var MaterialPalenightTheme = Theme{
	Name:                "Material Palenight",
	Handle:              "material-palenight",
	Background:          lipgloss.Color("#292D3E"),
	Foreground:          lipgloss.Color("#BFC7D5"),
	Accent:              lipgloss.Color("#82AAFF"),
	Secondary:           lipgloss.Color("#C792EA"),
	SelectionBackground: lipgloss.Color("#444267"),
}

var TokyoNightTheme = Theme{
	Name:                "Tokyo Night",
	Handle:              "tokyo-night",
	Background:          lipgloss.Color("#1A1B26"),
	Foreground:          lipgloss.Color("#C0CAF5"),
	Accent:              lipgloss.Color("#7AA2F7"),
	Secondary:           lipgloss.Color("#F7768E"),
	SelectionBackground: lipgloss.Color("#33467C"),
}

var Synthwave84Theme = Theme{
	Name:                "Synthwave '84",
	Handle:              "synthwave",
	Background:          lipgloss.Color("#2B213A"),
	Foreground:          lipgloss.Color("#F92AAD"),
	Accent:              lipgloss.Color("#FF6C7A"),
	Secondary:           lipgloss.Color("#F4F99D"),
	SelectionBackground: lipgloss.Color("#3E2A43"),
}

var Cobalt2Theme = Theme{
	Name:                "Cobalt2",
	Handle:              "cobalt2",
	Background:          lipgloss.Color("#193549"),
	Foreground:          lipgloss.Color("#FFFFFF"),
	Accent:              lipgloss.Color("#FF9D00"),
	Secondary:           lipgloss.Color("#FF0000"),
	SelectionBackground: lipgloss.Color("#1D3B53"),
}

var PaperColorDarkTheme = Theme{
	Name:                "PaperColor Dark",
	Handle:              "papercolor-dark",
	Background:          lipgloss.Color("#1C1C1C"),
	Foreground:          lipgloss.Color("#EEEEEE"),
	Accent:              lipgloss.Color("#87AF5F"),
	Secondary:           lipgloss.Color("#D7875F"),
	SelectionBackground: lipgloss.Color("#3A3A3A"),
}

var PaperColorLightTheme = Theme{
	Name:                "PaperColor Light",
	Handle:              "papercolor-light",
	Background:          lipgloss.Color("#EEEEEE"),
	Foreground:          lipgloss.Color("#1C1C1C"),
	Accent:              lipgloss.Color("#5F875F"),
	Secondary:           lipgloss.Color("#AF5F5F"),
	SelectionBackground: lipgloss.Color("#D7D7D7"),
}

var GitHubDarkTheme = Theme{
	Name:                "GitHub Dark",
	Handle:              "github-dark",
	Background:          lipgloss.Color("#0D1117"),
	Foreground:          lipgloss.Color("#C9D1D9"),
	Accent:              lipgloss.Color("#58A6FF"),
	Secondary:           lipgloss.Color("#F85149"),
	SelectionBackground: lipgloss.Color("#30363D"),
}
