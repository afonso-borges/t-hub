package themes

import (
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

var (
	Red     = lipgloss.Color("#D20A2E")
	NormalFg = lipgloss.AdaptiveColor{Light: "235", Dark: "252"}
	Indigo   = lipgloss.AdaptiveColor{Light: "#5A56E0", Dark: "#7571F9"}
	Fuchsia  = lipgloss.Color("#F780E2")
)

func DefaultTheme() *huh.Theme {
	t := *huh.ThemeCharm()

	t.Focused.SelectedPrefix = lipgloss.NewStyle().
		Foreground(Red).
		SetString("✗ ")

	t.Focused.SelectedOption = t.Focused.SelectedOption.
		Foreground(Red)

	t.Focused.UnselectedPrefix = lipgloss.NewStyle().
		Foreground(NormalFg).
		SetString("• ")

	t.Focused.UnselectedOption = t.Focused.UnselectedOption.
		Foreground(NormalFg)

	return &t
}