package zone

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"

	pb "github.com/trevatk/tbd/lib/protocol/dns/authoritative/v1"
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
	inputs := make([]textinput.Model, 2)
	inputs[domain] = textinput.New()
	inputs[domain].Focus()

	inputs[adminEmail] = textinput.New()

	return model{
		inputs:  inputs,
		focused: 0,
		err:     nil,
	}
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
			// if m.focused == len(m.inputs)-1 {
			// 	if client == nil {
			// 		// log error and quit
			// 	}

			// 	_, err := client.CreateZone(context.Background(), &pb.CreateZoneRequest{
			// 		DomainOrNamespace: m.inputs[domain].Value(),
			// 		AdminEmail:        m.inputs[adminEmail].Value(),
			// 	})
			// 	if err != nil {
			// 		return m.Update(err)
			// 	}
			// 	fmt.Println("zone created")
			// }
			// m.nextInput()
			if client == nil {
				// log error and quit
			}

			_, err := client.CreateZone(context.Background(), &pb.CreateZoneRequest{
				DomainOrNamespace: m.inputs[domain].Value(),
				AdminEmail:        m.inputs[adminEmail].Value(),
			})
			if err != nil {
				return m.Update(err)
			}

		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyTab:
			m.nextInput()
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

func (m *model) nextInput() {
	m.focused = (m.focused + 1) % len(m.inputs)
}

var (
	client pb.AuthoritativeServiceClient

	createCmd = &cobra.Command{
		Use: "create",
		RunE: func(cmd *cobra.Command, args []string) error {

			var err error
			client, err = newClient("localhost:50051")
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}

			p := tea.NewProgram(initialModel())
			if _, err := p.Run(); err != nil {
				return fmt.Errorf("%v", err)
			}
			return nil
		},
	}

	hintStyle  = lipgloss.NewStyle()
	inputStyle = lipgloss.NewStyle()
)
