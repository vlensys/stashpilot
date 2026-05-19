package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/vlensys/stashpilot/internal/git"
)

type appState int

const (
	stateBrowse appState = iota
	stateNewStash
	stateConfirmApply
	stateConfirmPop
	stateConfirmDrop
)

type (
	stashesLoadedMsg struct {
		stashes []git.Stash
		err     error
	}
	diffLoadedMsg struct {
		ref  string
		diff string
		err  error
	}
	actionDoneMsg struct {
		action string
		err    error
	}
)

type Model struct {
	stashes        []git.Stash
	cursor         int
	state          appState
	input          textinput.Model
	viewport       viewport.Model
	width          int
	height         int
	statusMsg      string
	statusErr      bool
	loadingDiff    bool
	currentDiffRef string
	ready          bool
}

func New() Model {
	ti := textinput.New()
	ti.Placeholder = "optional message…"
	ti.CharLimit = 100
	ti.Width = 40
	ti.PromptStyle = inputLabelStyle
	ti.Prompt = "> "

	return Model{
		input: ti,
		state: stateBrowse,
	}
}

func cmdLoadStashes() tea.Cmd {
	return func() tea.Msg {
		stashes, err := git.List()
		return stashesLoadedMsg{stashes: stashes, err: err}
	}
}

func cmdLoadDiff(ref string) tea.Cmd {
	return func() tea.Msg {
		diff, err := git.Diff(ref)
		return diffLoadedMsg{ref: ref, diff: diff, err: err}
	}
}

func cmdApply(ref string) tea.Cmd {
	return func() tea.Msg {
		err := git.Apply(ref)
		return actionDoneMsg{action: "apply", err: err}
	}
}

func cmdPop(index int) tea.Cmd {
	return func() tea.Msg {
		err := git.Pop(index)
		return actionDoneMsg{action: "pop", err: err}
	}
}

func cmdDrop(index int) tea.Cmd {
	return func() tea.Msg {
		err := git.Drop(index)
		return actionDoneMsg{action: "drop", err: err}
	}
}

func cmdPush(message string) tea.Cmd {
	return func() tea.Msg {
		err := git.Push(message)
		return actionDoneMsg{action: "push", err: err}
	}
}

func (m Model) Init() tea.Cmd {
	return cmdLoadStashes()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		pw := m.width - listWidth - 5
		ph := m.height - 4
		if !m.ready {
			m.viewport = viewport.New(pw, ph)
			m.viewport.Style = lipgloss.NewStyle().Padding(0, 1)
			m.ready = true
		} else {
			m.viewport.Width = pw
			m.viewport.Height = ph
		}

	case stashesLoadedMsg:
		if msg.err != nil {
			m.statusMsg = msg.err.Error()
			m.statusErr = true
		} else {
			prev := m.cursor
			m.stashes = msg.stashes
			if len(m.stashes) == 0 {
				m.cursor = 0
				m.viewport.SetContent(dimStyle.Render("no stashes to preview"))
			} else {
				if prev >= len(m.stashes) {
					m.cursor = len(m.stashes) - 1
				} else {
					m.cursor = prev
				}
				ref := m.stashes[m.cursor].Ref
				m.currentDiffRef = ref
				m.loadingDiff = true
				cmds = append(cmds, cmdLoadDiff(ref))
			}
		}

	case diffLoadedMsg:
		m.loadingDiff = false
		if msg.ref == m.currentDiffRef {
			if msg.err != nil {
				m.viewport.SetContent(statusErrStyle.Render("diff error: " + msg.err.Error()))
			} else {
				m.viewport.SetContent(colorDiff(msg.diff))
			}
			m.viewport.GotoTop()
		}

	case actionDoneMsg:
		m.state = stateBrowse
		if msg.err != nil {
			m.statusMsg = fmt.Sprintf("%s failed: %s", msg.action, msg.err.Error())
			m.statusErr = true
		} else {
			m.statusMsg = fmt.Sprintf("stash %sd successfully", msg.action)
			m.statusErr = false
		}
		cmds = append(cmds, cmdLoadStashes())

	case tea.KeyMsg:
		switch m.state {
		case stateNewStash:
			return m.updateInput(msg, cmds)
		case stateConfirmApply, stateConfirmPop, stateConfirmDrop:
			return m.updateConfirm(msg, cmds)
		default:
			return m.updateBrowse(msg, cmds)
		}
	}

	var vpCmd tea.Cmd
	m.viewport, vpCmd = m.viewport.Update(msg)
	cmds = append(cmds, vpCmd)

	return m, tea.Batch(cmds...)
}

func (m Model) updateBrowse(msg tea.KeyMsg, cmds []tea.Cmd) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, keys.Quit):
		return m, tea.Quit

	case key.Matches(msg, keys.Up):
		if m.cursor > 0 {
			m.cursor--
			m.statusMsg = ""
			m.loadingDiff = true
			ref := m.stashes[m.cursor].Ref
			m.currentDiffRef = ref
			cmds = append(cmds, cmdLoadDiff(ref))
		}

	case key.Matches(msg, keys.Down):
		if m.cursor < len(m.stashes)-1 {
			m.cursor++
			m.statusMsg = ""
			m.loadingDiff = true
			ref := m.stashes[m.cursor].Ref
			m.currentDiffRef = ref
			cmds = append(cmds, cmdLoadDiff(ref))
		}

	case key.Matches(msg, keys.Apply):
		if len(m.stashes) > 0 {
			m.state = stateConfirmApply
			m.statusMsg = ""
		}

	case key.Matches(msg, keys.Pop):
		if len(m.stashes) > 0 {
			m.state = stateConfirmPop
			m.statusMsg = ""
		}

	case key.Matches(msg, keys.Drop):
		if len(m.stashes) > 0 {
			m.state = stateConfirmDrop
			m.statusMsg = ""
		}

	case key.Matches(msg, keys.New):
		m.state = stateNewStash
		m.input.Reset()
		m.statusMsg = ""
		cmds = append(cmds, m.input.Focus())

	case key.Matches(msg, keys.Refresh):
		m.statusMsg = ""
		cmds = append(cmds, cmdLoadStashes())

	case key.Matches(msg, keys.ScrollUp):
		m.viewport.HalfViewUp()

	case key.Matches(msg, keys.ScrollDown):
		m.viewport.HalfViewDown()
	}

	return m, tea.Batch(cmds...)
}

func (m Model) updateConfirm(msg tea.KeyMsg, cmds []tea.Cmd) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y", "enter":
		if len(m.stashes) == 0 {
			m.state = stateBrowse
			break
		}
		s := m.stashes[m.cursor]
		switch m.state {
		case stateConfirmApply:
			cmds = append(cmds, cmdApply(s.Ref))
		case stateConfirmPop:
			cmds = append(cmds, cmdPop(s.Index))
		case stateConfirmDrop:
			cmds = append(cmds, cmdDrop(s.Index))
		}
		m.state = stateBrowse
	case "n", "esc", "q":
		m.state = stateBrowse
	}
	return m, tea.Batch(cmds...)
}

func (m Model) updateInput(msg tea.KeyMsg, cmds []tea.Cmd) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		message := strings.TrimSpace(m.input.Value())
		m.state = stateBrowse
		m.input.Blur()
		cmds = append(cmds, cmdPush(message))
	case "esc":
		m.state = stateBrowse
		m.input.Blur()
	default:
		var cmd tea.Cmd
		m.input, cmd = m.input.Update(msg)
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if !m.ready {
		return "\n  loading…"
	}

	header := m.renderHeader()
	body := m.renderBody()
	footer := m.renderFooter()

	return lipgloss.JoinVertical(lipgloss.Left, header, body, footer)
}

func (m Model) renderHeader() string {
	title := titleStyle.Render("stashpilot")
	sub := subtitleStyle.Render("  git stash manager")
	count := ""
	if len(m.stashes) > 0 {
		count = dimStyle.Render(fmt.Sprintf("  %d stash(es)", len(m.stashes)))
	}
	return lipgloss.JoinHorizontal(lipgloss.Center, title, sub, count)
}

func (m Model) renderBody() string {
	listPanel := m.renderList()
	previewPanel := m.renderPreview()
	return lipgloss.JoinHorizontal(lipgloss.Top, listPanel, previewPanel)
}

func (m Model) renderList() string {
	innerH := m.height - 4
	innerW := listWidth - 2

	startIdx := 0
	itemsPerPage := innerH
	if m.cursor >= itemsPerPage {
		startIdx = m.cursor - itemsPerPage + 1
	}

	var rows []string
	for i := startIdx; i < len(m.stashes) && len(rows) < itemsPerPage; i++ {
		rows = append(rows, m.renderItem(m.stashes[i], i == m.cursor, innerW))
	}

	if len(m.stashes) == 0 {
		rows = append(rows, dimStyle.Width(innerW).Render("no stashes yet — press n to create one"))
	}

	content := strings.Join(rows, "\n")
	for len(rows) < itemsPerPage {
		content += "\n"
		rows = append(rows, "")
	}

	return listBorderStyle.
		Width(listWidth).
		Height(innerH).
		Render(content)
}

func (m Model) renderItem(s git.Stash, selected bool, width int) string {
	var rStyle, bStyle, tStyle, mStyle lipgloss.Style
	if selected {
		rStyle = selectedRefStyle
		bStyle = selectedBranchStyle
		tStyle = selectedTimeStyle
		mStyle = selectedStyle
	} else {
		rStyle = refStyle
		bStyle = branchStyle
		tStyle = timeStyle
		mStyle = normalStyle
	}

	refStr := rStyle.Render(s.Ref)
	ageStr := tStyle.Render(formatAge(s.Date))
	topLine := refStr + "  " + ageStr

	msg := s.Message
	if s.Branch != "" {
		branchStr := bStyle.Render(s.Branch)
		rest := strings.TrimPrefix(s.Message, "WIP on "+s.Branch+":")
		rest = strings.TrimPrefix(rest, "On "+s.Branch+":")
		rest = strings.TrimSpace(rest)
		msg = branchStr + "  " + dimStyle.Render(rest)
	}

	item := topLine + "\n" + mStyle.Render(truncate(stripAnsi(msg), width-2))

	if selected {
		return selectedStyle.Width(width).Render(item)
	}
	return normalStyle.Width(width).Render(item)
}

func (m Model) renderPreview() string {
	pw := m.width - listWidth - 5
	ph := m.height - 4

	titleStr := "diff preview"
	if m.loadingDiff {
		titleStr = "loading diff…"
	} else if len(m.stashes) > 0 {
		titleStr = m.stashes[m.cursor].Ref
	}

	previewTitle := dimStyle.Render(titleStr)
	scroll := ""
	if m.viewport.TotalLineCount() > m.viewport.Height {
		pct := int(m.viewport.ScrollPercent() * 100)
		scroll = dimStyle.Render(fmt.Sprintf(" %d%%", pct))
	}

	headerLine := lipgloss.NewStyle().
		Width(pw - 2).
		Render(previewTitle + lipgloss.NewStyle().Width(pw-2-lipgloss.Width(previewTitle)-lipgloss.Width(scroll)).Render("") + scroll)

	content := headerLine + "\n" + m.viewport.View()

	return previewBorderStyle.
		Width(pw).
		Height(ph).
		Render(content)
}

func (m Model) renderFooter() string {
	var content string

	switch m.state {
	case stateNewStash:
		content = inputLabelStyle.Render("new stash") + "  " + m.input.View() + "  " + helpStyle.Render("enter↵ save · esc cancel")

	case stateConfirmApply:
		ref := ""
		if len(m.stashes) > 0 {
			ref = m.stashes[m.cursor].Ref
		}
		content = confirmStyle.Render("apply "+ref+"?") + "  " + helpStyle.Render("y/enter confirm · n/esc cancel")

	case stateConfirmPop:
		ref := ""
		if len(m.stashes) > 0 {
			ref = m.stashes[m.cursor].Ref
		}
		content = confirmStyle.Render("pop "+ref+"?") + "  " + helpStyle.Render("y/enter confirm · n/esc cancel")

	case stateConfirmDrop:
		ref := ""
		if len(m.stashes) > 0 {
			ref = m.stashes[m.cursor].Ref
		}
		content = confirmStyle.Render("drop "+ref+"? (irreversible)") + "  " + helpStyle.Render("y/enter confirm · n/esc cancel")

	default:
		var statusPart string
		if m.statusMsg != "" {
			if m.statusErr {
				statusPart = statusErrStyle.Render(m.statusMsg) + "  "
			} else {
				statusPart = statusOKStyle.Render(m.statusMsg) + "  "
			}
		}
		helpPart := helpStyle.Render("↑/↓") + " nav  " +
			helpStyle.Render("a") + " apply  " +
			helpStyle.Render("p") + " pop  " +
			helpStyle.Render("d") + " drop  " +
			helpStyle.Render("n") + " new  " +
			helpStyle.Render("r") + " refresh  " +
			helpStyle.Render("q") + " quit"
		content = statusPart + helpPart
	}

	return footerStyle.Width(m.width).Render(content)
}

func colorDiff(diff string) string {
	lines := strings.Split(diff, "\n")
	var sb strings.Builder

	for _, line := range lines {
		switch {
		case strings.HasPrefix(line, "+++") || strings.HasPrefix(line, "---"):
			sb.WriteString(diffFileStyle.Render(line))
		case strings.HasPrefix(line, "+"):
			sb.WriteString(diffAddStyle.Render(line))
		case strings.HasPrefix(line, "-"):
			sb.WriteString(diffRemoveStyle.Render(line))
		case strings.HasPrefix(line, "@@"):
			sb.WriteString(diffHunkStyle.Render(line))
		default:
			sb.WriteString(line)
		}
		sb.WriteByte('\n')
	}

	return sb.String()
}

func formatAge(t time.Time) string {
	if t.IsZero() {
		return "unknown"
	}
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		return fmt.Sprintf("%dm ago", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%dh ago", int(d.Hours()))
	case d < 7*24*time.Hour:
		return fmt.Sprintf("%dd ago", int(d.Hours()/24))
	case d < 30*24*time.Hour:
		return fmt.Sprintf("%dw ago", int(d.Hours()/(24*7)))
	default:
		return t.Format("Jan 2 2006")
	}
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
}

func stripAnsi(s string) string {
	var b strings.Builder
	inEsc := false
	for i := 0; i < len(s); i++ {
		if s[i] == '\x1b' {
			inEsc = true
			continue
		}
		if inEsc {
			if s[i] == 'm' {
				inEsc = false
			}
			continue
		}
		b.WriteByte(s[i])
	}
	return b.String()
}
