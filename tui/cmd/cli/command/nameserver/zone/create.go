package zone

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

type (
	errMsg error

	model struct {
		inputs  []textinput.Model
		focused int
		err     error
	}
)

const (
	domain = iota
	adminEmail
)

func initialModel() model {
	var m model
	inputs := make([]textinput.Model, 2)
	inputs[0] = textinput.New()
	inputs[1] = textinput.New()
	m.inputs = inputs
	m.focused = 0
	return m
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, len(m.inputs))
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.focused == len(m.inputs)-1 {
				return m, tea.Quit // submit form
			}
			m.nextInput()
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}

		for i := range m.inputs {
			m.inputs[i].Blur()
		}
		m.inputs[m.focused].Focus()

	case errMsg:
		m.err = msg
		return m, nil
	}

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	return fmt.Sprintf(`
	
	%s		%s
	--		--
	%s		%s
	`, "Domain", "Admin Email", m.inputs[domain].View(), m.inputs[adminEmail].View()) + "\n"
}

func (m model) nextInput() {
	m.focused = (m.focused + 1) % len(m.inputs)
}

var (
	createCmd = &cobra.Command{
		Use: "create",
		RunE: func(cmd *cobra.Command, args []string) error {
			p := tea.NewProgram(initialModel())
			if _, err := p.Run(); err != nil {
				return fmt.Errorf("%v", err)
			}
			return nil
		},
	}
)
