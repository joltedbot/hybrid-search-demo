package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/elastic/go-elasticsearch/v9"
)

var (
	titleStyle  = lipgloss.NewStyle().Bold(true).Margin(0, 1)
	helpStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Margin(1, 1)
	docStyle    = lipgloss.NewStyle().Margin(1, 2)
	separator   = lipgloss.NewStyle().Foreground(lipgloss.Color("238")).Render(strings.Repeat("─", 80))
	resultTitle = lipgloss.NewStyle().Bold(true)
	resultScore = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
)

type (
	searchResultMsg struct{ results Result }
	errorMsg        struct{ err error }
)

type model struct {
	esClient    *elasticsearch.Client
	searchIndex string
	textInput   textinput.Model
	spinner     spinner.Model
	viewport    viewport.Model
	isLoading   bool
	results     Result
	err         error
}

func initialModel(esClient *elasticsearch.Client, index string) model {

	ti := textinput.New()
	ti.Placeholder = "Enter search query..."
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 50

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return model{
		esClient:    esClient,
		searchIndex: index,
		textInput:   ti,
		spinner:     s,
		isLoading:   false,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			if m.textInput.Value() == "" {
				return m, nil
			}

			m.isLoading = true
			m.err = nil
			m.results = Result{} // Clear previous results

			return m, tea.Batch(
				m.spinner.Tick,
				runSearchCmd(m.esClient, m.searchIndex, m.textInput.Value()),
			)
		}
	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - 10 // Adjust for header and footer
		separator = lipgloss.NewStyle().Foreground(lipgloss.Color("238")).Render(strings.Repeat("─", msg.Width))
	case searchResultMsg:
		m.isLoading = false
		m.results = msg.results
		m.viewport.SetContent(m.formatResults())
		m.viewport.GotoTop()
		return m, nil
	case errorMsg:
		m.isLoading = false
		m.err = msg.err
		return m, nil
	case spinner.TickMsg:
		if m.isLoading {
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)
	cmds = append(cmds, cmd)

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {

	var b strings.Builder

	b.WriteString(titleStyle.Render("Hybrid Search Demo"))
	b.WriteString("\n")

	// Input
	if m.isLoading {
		b.WriteString(fmt.Sprintf("Search: %s %s", m.textInput.View(), m.spinner.View()))
	} else {
		b.WriteString(fmt.Sprintf("Search: %s", m.textInput.View()))
	}

	b.WriteString("\n\n")

	if m.err != nil {
		b.WriteString(fmt.Sprintf("Error: %v\n", m.err))
	} else {
		b.WriteString(m.viewport.View())
	}

	footerText := "Use ↑/↓ to scroll | Ctrl+C to quit"

	if len(m.results.Hits.Hits) > 0 {
		footerText = fmt.Sprintf("%d results | %s", len(m.results.Hits.Hits), footerText)
	}

	b.WriteString(helpStyle.Render(footerText))

	return b.String()
}

func (m model) formatResults() string {

	if len(m.results.Hits.Hits) == 0 && !m.isLoading {
		return "No results found."
	}

	var s strings.Builder

	for i, hit := range m.results.Hits.Hits {
		title, _ := hit.Source["Title"].(string)
		whatYouShouldDo, _ := hit.Source["What you should do"].(string)

		s.WriteString(docStyle.Render(
			resultTitle.Render(fmt.Sprintf("Result %d: %s", i+1, title)) + "\n" +
				resultScore.Render(fmt.Sprintf("Score: %.2f", hit.Score)),
		))

		if whatYouShouldDo != "" {
			s.WriteString(docStyle.Render(fmt.Sprintf("What to do: %s", whatYouShouldDo)))
		}

		if i < len(m.results.Hits.Hits)-1 {
			s.WriteString("\n" + separator + "\n")
		}
	}

	return s.String()
}

func runSearchCmd(esClient *elasticsearch.Client, index, query string) tea.Cmd {

	return func() tea.Msg {
		res, err := runQuery(esClient, index, query)
		if err != nil {
			return errorMsg{err}
		}

		var resultData Result

		if err := json.Unmarshal([]byte(res), &resultData); err != nil {
			return errorMsg{err}
		}

		return searchResultMsg{results: resultData}
	}
}
