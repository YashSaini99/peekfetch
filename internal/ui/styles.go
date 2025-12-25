package ui

import (
	"github.com/charmbracelet/lipgloss"
)

// Color palette - Calm, modern, professional colors
var (
	ColorPrimary   = lipgloss.Color("#7DCFFF") // Soft blue
	ColorSecondary = lipgloss.Color("#BB9AF7") // Soft purple
	ColorSuccess   = lipgloss.Color("#9ECE6A") // Soft green
	ColorWarning   = lipgloss.Color("#E0AF68") // Soft amber
	ColorDanger    = lipgloss.Color("#F7768E") // Soft red
	ColorInfo      = lipgloss.Color("#7AA2F7") // Medium blue
	ColorAccent    = lipgloss.Color("#73DACA") // Soft teal
	ColorMuted     = lipgloss.Color("#565F89") // Muted blue-gray
	ColorText      = lipgloss.Color("#C0CAF5") // Light periwinkle
	ColorSubtle    = lipgloss.Color("#9AA5CE") // Subtle text
	ColorBg        = lipgloss.Color("#1A1B26") // Deep blue-black
	ColorBorder    = lipgloss.Color("#414868") // Border gray
)

// Styles for the UI
var (
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary).
			PaddingLeft(1).
			PaddingRight(1)

	BannerStyle = lipgloss.NewStyle().
			Foreground(ColorPrimary).
			Bold(true).
			Align(lipgloss.Center)

	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			BorderStyle(lipgloss.DoubleBorder()).
			BorderForeground(ColorPrimary).
			Padding(0, 2).
			MarginBottom(1)

	SelectedStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorSecondary).
			Background(lipgloss.Color("#292E42")).
			PaddingLeft(1).
			PaddingRight(1).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderLeft(true).
			BorderForeground(ColorSecondary)

	NormalStyle = lipgloss.NewStyle().
			Foreground(ColorText).
			PaddingLeft(1)

	KeyStyle = lipgloss.NewStyle().
			Foreground(ColorAccent).
			Bold(true)

	ValueStyle = lipgloss.NewStyle().
			Foreground(ColorSuccess)

	TreeStyle = lipgloss.NewStyle().
			Foreground(ColorBorder)

	SubKeyStyle = lipgloss.NewStyle().
			Foreground(ColorInfo)

	FooterStyle = lipgloss.NewStyle().
			Foreground(ColorMuted).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderTop(true).
			BorderForeground(ColorPrimary).
			Padding(0, 2).
			MarginTop(1).
			Align(lipgloss.Center)

	LiveBadgeStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorBg).
			Background(ColorSuccess).
			Padding(0, 1).
			MarginLeft(1)

	SectionHeaderStyle = lipgloss.NewStyle().
				Foreground(ColorPrimary).
				Bold(true)

	ExpandedContentStyle = lipgloss.NewStyle().
				MarginLeft(2).
				BorderStyle(lipgloss.RoundedBorder()).
				BorderLeft(true).
				BorderForeground(ColorBorder).
				PaddingLeft(2).
				MarginTop(1).
				MarginBottom(1)

	ProgressBarStyle = lipgloss.NewStyle().
				Foreground(ColorSuccess)

	ProgressEmptyStyle = lipgloss.NewStyle().
				Foreground(ColorMuted)
)
