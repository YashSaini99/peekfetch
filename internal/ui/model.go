package ui

import (
	"time"

	"peekfetch/internal/sysinfo"
	"peekfetch/internal/types"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	Sections       []types.Section
	SelectedIndex  int
	ScrollOffset   int // Scroll offset for expanded content
	LiveMode       bool
	Width          int
	Height         int
	ViewportHeight int // Available height for content
}

type tickMsg time.Time

// TickCmd returns a command that sends a tick message every 500ms
func tickCmd() tea.Cmd {
	return tea.Tick(500*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// InitialModel creates the initial model with all sections
func InitialModel() Model {
	return Model{
		Sections: []types.Section{
			sysinfo.GetSystemInfo(),
			sysinfo.GetCPUInfo(),
			sysinfo.GetMemoryInfo(),
			sysinfo.GetDiskInfo(),
			sysinfo.GetNetworkInfo(),
		},
		SelectedIndex:  0,
		ScrollOffset:   0,
		LiveMode:       false,
		Width:          80,
		Height:         24,
		ViewportHeight: 20,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}
