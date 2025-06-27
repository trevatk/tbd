package server

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/codes"
)

const (
	url = "https://google.com"
)

type (
	statusMsg uint32

	errMsg struct{ error }

	pingModel struct {
		ctx context.Context

		statusCode codes.Code
		err        error
	}
)

func (m pingModel) Init() tea.Cmd {
	return checkServerWithContext(m.ctx)
}

func (m pingModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		default:
			return m, nil
		}

	case statusMsg:
		m.statusCode = codes.Code(msg)
		return m, tea.Quit

	case errMsg:
		m.err = msg
		return m, nil

	default:
		return m, nil
	}
}

func (m pingModel) View() string {
	s := fmt.Sprintf("check url %s...", url)
	if m.err != nil {
		s += fmt.Sprintf("something went wrong %s", m.err)
	} else if m.statusCode != codes.OK {
		s += fmt.Sprintf("%d %s", m.statusCode, m.statusCode.String())
	}
	return s + "\n"
}

func checkServerWithContext(_ context.Context) tea.Cmd {
	return func() tea.Msg {
		return statusMsg(codes.OK)
	}
}

var (
	pingCmd = &cobra.Command{
		Use: "ping",
		RunE: func(cmd *cobra.Command, args []string) error {
			p := tea.NewProgram(pingModel{ctx: cmd.Context()})
			_, err := p.Run()
			if err != nil {
				return fmt.Errorf("failed to create tea program: %w", err)
			}
			return nil
		},
	}
)
