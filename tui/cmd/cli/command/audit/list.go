package audit

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/codes"
)

var (
	limit  int64
	offset int64
)

type (
	listStatusMsg uint32

	listModel struct {
		ctx        context.Context
		serverAddr string

		statusCode codes.Code
		err        error
	}

	errMsg struct{ error }
)

func (m listModel) Init() tea.Cmd {
	return listTxs(m.ctx, m.serverAddr, limit, offset)
}

func (m listModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		default:
			return m, nil
		}
	case listStatusMsg:
		m.statusCode = codes.Code(msg)
		return m, tea.Quit
	case errMsg:
		m.err = msg
		return m, nil

	default:
		return m, nil
	}
}

func (m listModel) View() string {
	s := "list transactions..."
	if m.err != nil {
		s += fmt.Sprintf("error occured %s", m.err)
	} else if m.statusCode != codes.OK {
		s += fmt.Sprintf("%d %s", m.statusCode, codes.Code(m.statusCode).String())
	}
	return s
}

func listTxs(_ context.Context, _ string, _, _ int64) tea.Cmd {
	return func() tea.Msg {
		// client, err := audit.NewClient(addr)
		// if err != nil {
		// 	return errMsg{err}
		// }

		// txs, err := client.List(ctx, limit, offset)
		// if err != nil {
		// 	return errMsg{err}
		// }
		// fmt.Println(txs)
		return listStatusMsg(codes.OK)
	}
}

var listCmd = &cobra.Command{
	Use: "list",
	RunE: func(cmd *cobra.Command, args []string) error {
		p := tea.NewProgram(listModel{ctx: cmd.Context(), serverAddr: "localhost:8000"})
		_, err := p.Run()
		if err != nil {
			return fmt.Errorf("failed to run program: %w", err)
		}
		return nil
	},
}
