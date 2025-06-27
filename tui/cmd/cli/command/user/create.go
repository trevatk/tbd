package user

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"google.golang.org/grpc/codes"

	"github.com/structx/tbd/lib/protocol"
	pb "github.com/structx/tbd/lib/protocol/identities/v1"
)

var (
	userEmail string
)

type (
	createStatusMsg int

	createUserModel struct {
		ctx        context.Context
		serverAddr string
		email      string

		statusCode int
		err        error
	}

	errMsg struct{ error }
)

func (m createUserModel) Init() tea.Cmd {
	return createUser(m.ctx, m.serverAddr, m.email)
}

func (m createUserModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		default:
			return m, nil
		}
	case createStatusMsg:
		m.statusCode = int(msg)
		return m, tea.Quit
	case errMsg:
		m.err = msg
		return m, nil

	default:
		return m, nil
	}
}

func (m createUserModel) View() string {
	s := fmt.Sprintf("create user %s...", m.email)
	if m.err != nil {
		s += fmt.Sprintf("error occured %s", m.err)
	} else if m.statusCode != 0 {
		s += fmt.Sprintf("%d %s", m.statusCode, codes.Code(m.statusCode).String())
	}
	return s
}

func createUser(ctx context.Context, addr, email string) tea.Cmd {
	return func() tea.Msg {
		conn, err := protocol.NewConn(addr)
		if err != nil {
			return errMsg{err}
		}

		client := pb.NewIdentitiesServiceClient(conn)
		_ = client
		// _, err = client.CreateUser(ctx, email)
		// if err != nil {
		// return errMsg{err}
		// }

		return createStatusMsg(codes.OK)
	}
}

var (
	createCmd = &cobra.Command{
		Use: "create",
		// Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			m := createUserModel{
				ctx:        cmd.Context(),
				serverAddr: "localhost:8000",
				email:      userEmail,
			}
			p := tea.NewProgram(m)
			_, err := p.Run()
			if err != nil {
				return fmt.Errorf("failed to execute create user command: %w", err)
			}
			return nil
		},
	}
)
