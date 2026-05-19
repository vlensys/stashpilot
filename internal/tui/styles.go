package tui

import "github.com/charmbracelet/lipgloss"

const listWidth = 44

var (
	subtle    = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#3D3D3D"}
	highlight = lipgloss.Color("#7D56F4")

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(highlight).
			Padding(0, 2)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#696969", Dark: "#909090"})

	selectedStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(highlight)

	normalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#1A1A1A", Dark: "#DDDDDD"})

	dimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#9B9B9B", Dark: "#5C5C5C"})

	refStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4")).
			Bold(true)

	selectedRefStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#C4B5FD")).
				Bold(true)

	branchStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#73F59F"))

	selectedBranchStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#A7F3D0"))

	timeStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#9B9B9B", Dark: "#737373"})

	selectedTimeStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#DDD6FE"))

	listBorderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder(), true, true, true, true).
			BorderForeground(subtle)

	previewBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder(), true, true, true, false).
				BorderForeground(subtle)

	footerStyle = lipgloss.NewStyle().
			Background(lipgloss.AdaptiveColor{Light: "#E8E8E8", Dark: "#2A2A2A"}).
			Foreground(lipgloss.AdaptiveColor{Light: "#444444", Dark: "#AAAAAA"}).
			Padding(0, 1)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#9B9B9B", Dark: "#606060"})

	keyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#555555", Dark: "#B0B0B0"}).
			Bold(true)

	statusOKStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#73F59F")).
			Bold(true)

	statusErrStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF6B6B")).
			Bold(true)

	confirmStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFD700")).
			Bold(true)

	inputLabelStyle = lipgloss.NewStyle().
			Foreground(highlight).
			Bold(true)

	diffAddStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#73F59F"))

	diffRemoveStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF6B6B"))

	diffHunkStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#89CFF0")).
			Bold(true)

	diffFileStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFD700")).
			Bold(true)
)
