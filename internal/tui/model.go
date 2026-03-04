package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/vertoo/skill-link/internal/core"
)

const minWidthForPreview = 80

var (
	titleStyle       = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	cursorStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
	symlinkStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))  // Green
	copyStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("39"))  // Blue
	mismatchStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("214")) // Orange/Yellow
	uninstalledStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240")) // Gray
	panelBorder      = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("240"))
	previewTitle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("99"))
	dimStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
)

type viewState int

const (
	viewList viewState = iota
	viewActionMenu
)

type actionType int

const (
	actionInstallSymlink actionType = iota
	actionInstallCopy
	actionUpdateCopy
	actionSwitchToSymlink
	actionSwitchToCopy
	actionUninstall
	actionCancel
)

type actionItem struct {
	label string
	act   actionType
}

type model struct {
	skills       []SkillItem
	cursor       int
	filter       string
	state        viewState
	selectedItem *SkillItem
	menuCursor   int
	actions      []actionItem
	err          error
	message      string
	width        int
	height       int
}

func initialModel() model {
	skills, err := loadSkills()
	return model{
		skills: skills,
		state:  viewList,
		err:    err,
		width:  80,
		height: 24,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			if m.state == viewActionMenu {
				m.state = viewList
				m.message = ""
				return m, nil
			}
			return m, tea.Quit
		}
	}

	if m.state == viewList {
		return updateList(m, msg)
	} else if m.state == viewActionMenu {
		return updateMenu(m, msg)
	}

	return m, nil
}

func updateList(m model, msg tea.Msg) (tea.Model, tea.Cmd) {
	filtered := getFilteredSkills(m.skills, m.filter)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(filtered)-1 {
				m.cursor++
			}
		case "enter":
			if len(filtered) > 0 {
				m.selectedItem = &filtered[m.cursor]
				m.actions = buildActions(m.selectedItem.Status)
				m.menuCursor = 0
				m.state = viewActionMenu
				m.message = ""
			}
		case "backspace":
			if len(m.filter) > 0 {
				m.filter = m.filter[:len(m.filter)-1]
				m.cursor = 0
			}
		default:
			if len(msg.String()) == 1 {
				m.filter += msg.String()
				m.cursor = 0
			}
		}
	}
	return m, nil
}

func updateMenu(m model, msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.menuCursor > 0 {
				m.menuCursor--
			}
		case "down", "j":
			if m.menuCursor < len(m.actions)-1 {
				m.menuCursor++
			}
		case "enter":
			act := m.actions[m.menuCursor]
			if act.act == actionCancel {
				m.state = viewList
				return m, nil
			}
			err := executeAction(m.selectedItem.Name, act.act)
			if err != nil {
				m.message = "Error: " + err.Error()
			} else {
				m.message = "Success: " + m.selectedItem.Name
			}
			m.skills, _ = loadSkills()
			m.state = viewList
		}
	}
	return m, nil
}

func getFilteredSkills(skills []SkillItem, filter string) []SkillItem {
	if filter == "" {
		return skills
	}
	var filtered []SkillItem
	for _, s := range skills {
		if strings.Contains(strings.ToLower(s.Name), strings.ToLower(filter)) {
			filtered = append(filtered, s)
		}
	}
	return filtered
}

func buildActions(status SkillStatus) []actionItem {
	switch status {
	case StatusNotInstalled:
		return []actionItem{
			{"Install as Symlink (Recommended)", actionInstallSymlink},
			{"Install as Copy (Isolated)", actionInstallCopy},
			{"Cancel", actionCancel},
		}
	case StatusSymlink:
		return []actionItem{
			{"Change to physical Copy", actionSwitchToCopy},
			{"Uninstall", actionUninstall},
			{"Cancel", actionCancel},
		}
	case StatusCopyMatch:
		return []actionItem{
			{"Change to Symlink", actionSwitchToSymlink},
			{"Uninstall", actionUninstall},
			{"Cancel", actionCancel},
		}
	case StatusCopyMismatch:
		return []actionItem{
			{"Update (Overwrite local changes)", actionUpdateCopy},
			{"Change to Symlink (Overwrite local changes)", actionSwitchToSymlink},
			{"Uninstall (local copy has modifications!)", actionUninstall},
			{"Cancel", actionCancel},
		}
	}
	return []actionItem{{"Cancel", actionCancel}}
}

// renderLeftPanel renders the skill list or action menu
func (m model) renderLeftPanel(panelWidth int) string {
	var sb strings.Builder

	if m.state == viewList {
		sb.WriteString(titleStyle.Render("Skill-Link: Manage Skills"))
		sb.WriteString("\n")

		if m.message != "" {
			sb.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Render(m.message) + "\n")
		}

		sb.WriteString(fmt.Sprintf("Search: %s\n\n", m.filter))

		filtered := getFilteredSkills(m.skills, m.filter)
		if len(filtered) == 0 {
			sb.WriteString("  No results found.\n")
		}

		for i, s := range filtered {
			cursor := "  "
			if m.cursor == i {
				cursor = cursorStyle.Render("> ")
			}

			statusStr := ""
			switch s.Status {
			case StatusNotInstalled:
				statusStr = uninstalledStyle.Render("[ ]")
			case StatusSymlink:
				statusStr = symlinkStyle.Render("[S]")
			case StatusCopyMatch:
				statusStr = copyStyle.Render("[C]")
			case StatusCopyMismatch:
				statusStr = mismatchStyle.Render("[!]")
			}

			sb.WriteString(fmt.Sprintf("%s%s %s\n", cursor, statusStr, s.Name))
		}

		sb.WriteString("\n")
		sb.WriteString(uninstalledStyle.Render("[ ]") + " uninstalled  ")
		sb.WriteString(symlinkStyle.Render("[S]") + " symlink  ")
		sb.WriteString(copyStyle.Render("[C]") + " copy  ")
		sb.WriteString(mismatchStyle.Render("[!]") + " modified\n")
		sb.WriteString("[j/k] navigate • [enter] options • [type to search] • [esc] quit")
	} else if m.state == viewActionMenu {
		sb.WriteString(titleStyle.Render(fmt.Sprintf("Options for: %s", m.selectedItem.Name)))
		sb.WriteString("\n\n")

		for i, act := range m.actions {
			cursor := "  "
			if m.menuCursor == i {
				cursor = cursorStyle.Render("> ")
			}
			sb.WriteString(fmt.Sprintf("%s%s\n", cursor, act.label))
		}
		sb.WriteString("\n[j/k] navigate • [enter] select • [esc] back")
	}

	return sb.String()
}

// renderRightPanel renders the SKILL.md preview for the currently highlighted skill
func (m model) renderRightPanel(panelWidth, panelHeight int) string {
	filtered := getFilteredSkills(m.skills, m.filter)
	if len(filtered) == 0 {
		return dimStyle.Render("No skill selected")
	}

	var skillName string
	if m.state == viewActionMenu && m.selectedItem != nil {
		skillName = m.selectedItem.Name
	} else {
		skillName = filtered[m.cursor].Name
	}

	globalSkillsDir, _ := core.GetGlobalSkillsDir()
	skillMdPath := filepath.Join(globalSkillsDir, skillName, "SKILL.md")

	content, err := os.ReadFile(skillMdPath)
	if err != nil {
		return dimStyle.Render("No SKILL.md found")
	}

	// Render markdown with glamour (use fixed dark style to avoid terminal query conflicts)
	renderer, err := glamour.NewTermRenderer(
		glamour.WithStylePath("dark"),
		glamour.WithWordWrap(panelWidth-4),
	)
	if err != nil {
		return string(content)
	}

	rendered, err := renderer.Render(string(content))
	if err != nil {
		return string(content)
	}

	// Truncate to fit panel height
	lines := strings.Split(rendered, "\n")
	maxLines := panelHeight - 2
	if maxLines < 1 {
		maxLines = 1
	}
	if len(lines) > maxLines {
		lines = lines[:maxLines]
		lines = append(lines, dimStyle.Render("  ···"))
	}

	return strings.Join(lines, "\n")
}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Initialization error: %v\n", m.err)
	}

	showPreview := m.width >= minWidthForPreview

	if !showPreview {
		// Narrow terminal: single panel with padding
		content := m.renderLeftPanel(m.width - 4)
		return lipgloss.NewStyle().Padding(1, 2).Render(content)
	}

	// Wide terminal: split panels
	leftWidth := m.width*2/5 - 4
	rightWidth := m.width*3/5 - 4
	panelHeight := m.height - 4

	leftContent := m.renderLeftPanel(leftWidth)
	rightContent := m.renderRightPanel(rightWidth, panelHeight)

	leftPanel := panelBorder.
		Width(leftWidth).
		Height(panelHeight).
		Padding(1, 1).
		Render(leftContent)

	rightPanel := panelBorder.
		Width(rightWidth).
		Height(panelHeight).
		Padding(1, 1).
		Render(rightContent)

	return lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, rightPanel)
}
