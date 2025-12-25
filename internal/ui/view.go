package ui

import (
	"fmt"
	"strconv"
	"strings"

	"peekfetch/internal/sysinfo"
	"peekfetch/internal/types"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		// Reserve space for header (6 lines) and footer (3 lines)
		m.ViewportHeight = msg.Height - 9
		if m.ViewportHeight < 5 {
			m.ViewportHeight = 5
		}
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("q", "Q", "ctrl+c"))):
			return m, tea.Quit

		case key.Matches(msg, key.NewBinding(key.WithKeys("up", "k"))):
			// If section is expanded and has scroll, scroll up first
			if m.Sections[m.SelectedIndex].Expanded && m.ScrollOffset > 0 {
				m.ScrollOffset--
			} else if m.SelectedIndex > 0 {
				// Move to previous section
				m.SelectedIndex--
				m.ScrollOffset = 0 // Reset scroll for new section
			}

		case key.Matches(msg, key.NewBinding(key.WithKeys("down", "j"))):
			// If section is expanded, check if we can scroll down
			if m.Sections[m.SelectedIndex].Expanded {
				contentLines := m.countContentLines(m.Sections[m.SelectedIndex])
				if contentLines > m.ViewportHeight && m.ScrollOffset < contentLines-m.ViewportHeight {
					m.ScrollOffset++
					return m, nil
				}
			}
			// Move to next section
			if m.SelectedIndex < len(m.Sections)-1 {
				m.SelectedIndex++
				m.ScrollOffset = 0 // Reset scroll for new section
			}

		case key.Matches(msg, key.NewBinding(key.WithKeys("pageup", "ctrl+u"))):
			// Page up scrolling
			if m.Sections[m.SelectedIndex].Expanded && m.ScrollOffset > 0 {
				m.ScrollOffset -= m.ViewportHeight / 2
				if m.ScrollOffset < 0 {
					m.ScrollOffset = 0
				}
			}

		case key.Matches(msg, key.NewBinding(key.WithKeys("pagedown", "ctrl+d"))):
			// Page down scrolling
			if m.Sections[m.SelectedIndex].Expanded {
				contentLines := m.countContentLines(m.Sections[m.SelectedIndex])
				maxScroll := contentLines - m.ViewportHeight
				if maxScroll > 0 {
					m.ScrollOffset += m.ViewportHeight / 2
					if m.ScrollOffset > maxScroll {
						m.ScrollOffset = maxScroll
					}
				}
			}

		case key.Matches(msg, key.NewBinding(key.WithKeys("enter", " "))):
			m.Sections[m.SelectedIndex].Expanded = !m.Sections[m.SelectedIndex].Expanded
			m.ScrollOffset = 0 // Reset scroll when toggling

		case key.Matches(msg, key.NewBinding(key.WithKeys("l", "L"))):
			m.LiveMode = !m.LiveMode
			if m.LiveMode {
				return m, tickCmd()
			}
		}

	case tickMsg:
		if m.LiveMode {
			// Update CPU section
			for i := range m.Sections {
				if m.Sections[i].Name == "CPU" && m.Sections[i].LiveData {
					if m.Sections[i].Data == nil {
						m.Sections[i].Data = make(map[string]string)
					}
					m.Sections[i].Data["Usage"] = sysinfo.GetCPUUsage()
				}
				if m.Sections[i].Name == "Memory" && m.Sections[i].LiveData {
					updatedSection := sysinfo.GetMemoryInfo()
					m.Sections[i].Data = updatedSection.Data
				}
			}
			return m, tickCmd()
		}
	}

	return m, nil
}

func (m Model) View() string {
	var b strings.Builder

	// ASCII Banner
	banner := `
   ___           _   ___     _       _     
  / _ \___ ___ | | _| __|__| |_ ___| |__  
 | (_) / -_) -_)| |/ / _|/ -_)  _/ _| '_ \ 
  \___/\___\___||___/|___\___\__\__|_| |_|
`
	b.WriteString(BannerStyle.Render(banner))
	b.WriteString("\n")

	// Header with live badge
	title := "⚡ Interactive System Information"
	if m.LiveMode {
		title += LiveBadgeStyle.Render("● LIVE")
	}
	header := HeaderStyle.Render(title)
	b.WriteString(header)
	b.WriteString("\n")

	// Sections
	for i, section := range m.Sections {
		isSelected := i == m.SelectedIndex

		// Section name with icon
		icon := getSectionIcon(section.Name)
		expandIndicator := "▸"
		if section.Expanded {
			expandIndicator = "▾"
		}

		sectionLine := fmt.Sprintf(" %s  %s %s", expandIndicator, icon, section.Name)

		if isSelected {
			sectionLine = SelectedStyle.Render(sectionLine)
		} else {
			sectionLine = NormalStyle.Render(sectionLine)
		}

		b.WriteString(sectionLine)
		b.WriteString("\n")

		// Section content (if expanded)
		if section.Expanded && isSelected {
			var content strings.Builder

			// Get all content lines
			allLines := m.renderSectionContent(section)

			// Apply scrolling - only show visible lines
			startLine := m.ScrollOffset
			endLine := m.ScrollOffset + m.ViewportHeight
			if endLine > len(allLines) {
				endLine = len(allLines)
			}
			if startLine < len(allLines) {
				visibleLines := allLines[startLine:endLine]
				for _, line := range visibleLines {
					content.WriteString(line)
					content.WriteString("\n")
				}

				// Add scroll indicators
				if m.ScrollOffset > 0 {
					indicator := lipgloss.NewStyle().Foreground(ColorMuted).Render("    ▲ More above (↑ or PgUp)")
					content.WriteString(indicator + "\n")
				}
				if endLine < len(allLines) {
					indicator := lipgloss.NewStyle().Foreground(ColorMuted).Render("    ▼ More below (↓ or PgDn)")
					content.WriteString(indicator + "\n")
				}
			}

			b.WriteString(ExpandedContentStyle.Render(content.String()))
			b.WriteString("\n")
		}
	}

	// Footer
	footerText := "↑↓ Navigate/Scroll  │  PgUp/PgDn Fast Scroll  │  ⏎ Expand  │  L Live  │  Q Quit"
	footer := FooterStyle.Render(footerText)
	b.WriteString(footer)

	return b.String()
}

// getSectionIcon returns an icon for each section
func getSectionIcon(name string) string {
	icons := map[string]string{
		"System":  "🖥️ ",
		"CPU":     "⚡",
		"Memory":  "💾",
		"Disk":    "💿",
		"Network": "🌐",
	}
	if icon, ok := icons[name]; ok {
		return icon
	}
	return "📊"
}

// createProgressBar creates a visual progress bar
func createProgressBar(percent float64, width int) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}

	filledWidth := int((percent / 100.0) * float64(width))
	emptyWidth := width - filledWidth

	var color string
	if percent >= 90 {
		color = "#F7768E" // Soft red (danger)
	} else if percent >= 75 {
		color = "#E0AF68" // Soft amber (warning)
	} else if percent >= 50 {
		color = "#7DCFFF" // Soft blue (info)
	} else {
		color = "#9ECE6A" // Soft green (success)
	}

	filled := strings.Repeat("█", filledWidth)
	empty := strings.Repeat("░", emptyWidth)

	barStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(color))
	emptyStyle := lipgloss.NewStyle().Foreground(ColorMuted)

	return "[" + barStyle.Render(filled) + emptyStyle.Render(empty) + "]"
}

// countContentLines counts how many lines an expanded section would have
func (m Model) countContentLines(section types.Section) int {
	lines := 0

	if section.UseTree {
		for _, item := range section.TreeData {
			lines++                  // Item name
			lines += len(item.Order) // Children
		}
	} else {
		lines = len(section.Order)
		if lines == 0 {
			lines = len(section.Data)
		}
	}

	return lines
}

// renderSectionContent renders the section content and returns all lines
func (m Model) renderSectionContent(section types.Section) []string {
	allLines := []string{}

	if section.UseTree {
		// Tree structure rendering
		for i, item := range section.TreeData {
			isLast := i == len(section.TreeData)-1

			// Item name
			treeBranch := "├─"
			if isLast {
				treeBranch = "└─"
			}
			allLines = append(allLines, TreeStyle.Render(treeBranch)+" "+KeyStyle.Render(item.Name))

			// Children
			maxKeyLen := 0
			for key := range item.Children {
				if len(key) > maxKeyLen {
					maxKeyLen = len(key)
				}
			}

			for j, key := range item.Order {
				value := item.Children[key]
				isLastChild := j == len(item.Order)-1

				childBranch := "│  ├─"
				if isLast {
					childBranch = "   ├─"
				}
				if isLastChild {
					if isLast {
						childBranch = "   └─"
					} else {
						childBranch = "│  └─"
					}
				}

				padding := strings.Repeat(" ", maxKeyLen-len(key))

				// Add progress bar for percentage values
				if strings.HasSuffix(value, "%") && strings.Contains(key, "Usage") {
					percentStr := strings.TrimSuffix(value, "%")
					if percent, err := strconv.ParseFloat(percentStr, 64); err == nil {
						progressBar := createProgressBar(percent, 18)
						line := fmt.Sprintf("%s %s%s %s %s",
							TreeStyle.Render(childBranch),
							SubKeyStyle.Render(key),
							padding,
							progressBar,
							ValueStyle.Render(value))
						allLines = append(allLines, line)
						continue
					}
				}

				line := fmt.Sprintf("%s %s%s %s %s",
					TreeStyle.Render(childBranch),
					SubKeyStyle.Render(key),
					padding,
					TreeStyle.Render("│"),
					ValueStyle.Render(value))
				allLines = append(allLines, line)
			}
		}
	} else {
		// Regular key-value rendering
		maxKeyLen := 0
		for key := range section.Data {
			if len(key) > maxKeyLen {
				maxKeyLen = len(key)
			}
		}

		keys := section.Order
		if len(keys) == 0 {
			keys = []string{}
			for key := range section.Data {
				keys = append(keys, key)
			}
		}

		for _, key := range keys {
			value, ok := section.Data[key]
			if !ok {
				continue
			}

			padding := strings.Repeat(" ", maxKeyLen-len(key))

			// Add progress bar for percentage values
			if strings.HasSuffix(value, "%") && strings.Contains(key, "Usage") {
				percentStr := strings.TrimSuffix(value, "%")
				if percent, err := strconv.ParseFloat(percentStr, 64); err == nil {
					progressBar := createProgressBar(percent, 20)
					line := fmt.Sprintf("%s%s %s %s",
						KeyStyle.Render(key),
						padding,
						progressBar,
						ValueStyle.Render(value))
					allLines = append(allLines, line)
					continue
				}
			}

			line := fmt.Sprintf("%s%s %s %s",
				KeyStyle.Render(key),
				padding,
				KeyStyle.Render("│"),
				ValueStyle.Render(value))
			allLines = append(allLines, line)
		}
	}

	return allLines
}
