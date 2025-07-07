package chat

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"

	"github.com/trevatk/tbd/lib/protocol/chat"

	"github.com/structx/tbd/tui/cmd/cli/command"
	"github.com/structx/tbd/tui/internal/pkg/logging"
)

type sessionState uint
type item string

func (i item) FilterValue() string { return "" }

const (
	initialViews = 1

	sidebarView sessionState = iota
	threadView

	gap = "\n\n"
)

var (
	// ChatCmd chase cobra cli command
	ChatCmd = &cobra.Command{
		Use: "chat",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			client, err := chat.NewClient("did:grpc", "")
			if err != nil {
				return fmt.Errorf("failed to create chat cient: %w", err)
			}
			ctx = context.WithValue(ctx, "chatClient", client)
			cmd.SetContext(ctx)
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			cc, ok := ctx.Value("chatClient").(chat.Client)
			if !ok {
				return errors.New("chat client not set")
			}

			eventCh, err := cc.Subscribe(ctx)
			if err != nil {
				return fmt.Errorf("failed to create subscription: %w", err)
			}
			_ = eventCh

			p := tea.NewProgram(newMainModel(), tea.WithAltScreen())
			if _, err := p.Run(); err != nil {
				logging.FromContext(ctx).ErrorContext(ctx, "tea program", slog.String("error", err.Error()))
			}
			return nil
		},
	}

	modelStyle = lipgloss.NewStyle().
			Width(15).
			Height(5).
			Align(lipgloss.Center, lipgloss.Center).
			BorderStyle(lipgloss.HiddenBorder())
	focusedModelStyle = lipgloss.NewStyle().
				Width(15).
				Height(5).
				Align(lipgloss.Center, lipgloss.Center).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("69"))
	helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))

	docStyle = lipgloss.NewStyle().Margin(1, 2)
)

type (
	errMsg error

	mainModel struct {
		sub       chan struct{}
		responses int

		state sessionState

		sidebar sidebarModel
		index   int

		textarea textarea.Model
		messages []string

		viewport viewport.Model

		err errMsg
	}

	sidebarModel struct {
		list   list.Model
		active string
		quit   bool
	}
)

func (m sidebarModel) Init() tea.Cmd {
	return nil
}

func newMainModel() mainModel {
	m := mainModel{state: sidebarView, err: nil}
	m.index = 0
	m.sub = make(chan struct{})
	m.responses = 0
	m.sidebar = newSidebarModel()

	m.textarea = textarea.New()
	m.textarea.Placeholder = "Send a message..."
	// m.textarea.Focus()
	m.textarea.Prompt = "|"
	m.textarea.CharLimit = 280
	m.textarea.SetWidth(30)
	m.textarea.SetHeight(3)
	m.textarea.FocusedStyle.CursorLine = lipgloss.NewStyle()
	m.textarea.ShowLineNumbers = false
	m.textarea.KeyMap.InsertNewline.SetEnabled(false)

	m.viewport = viewport.New(30, 5)
	return m
}

func newSidebarModel() sidebarModel {
	m := sidebarModel{}
	items := []list.Item{
		item("Ramen"),
		item("Tomato Soup"),
		item("Hamburgers"),
		item("Cheeseburgers"),
		item("Currywurst"),
		item("Okonomiyaki"),
		item("Pasta"),
		item("Fillet Mignon"),
		item("Caviar"),
		item("Just Wine"),
	}

	m.list = list.New(items, list.NewDefaultDelegate(), 14, 20)
	m.list.Title = "threads"
	m.list.SetShowStatusBar(true)
	return m
}

func (m mainModel) Init() tea.Cmd {
	return tea.Batch(
		m.sidebar.Init(),
		waitForActivity(m.sub),
	)
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	// var cmd tea.Cmd
	var cmds []tea.Cmd
	cmds = append(cmds, tiCmd, vpCmd)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		m.textarea.SetWidth(msg.Width)
		m.viewport.Height = msg.Height - m.textarea.Height() - lipgloss.Height(gap)
		m.viewport.GotoBottom()
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "tab":
			if m.state == sidebarView {
				m.state = threadView
				m.textarea.Focus()
			} else {
				m.state = sidebarView
			}
		case "n":
			// if m.state == sidebarView {
			// 	cmds = append(cmds, m.sidebar.Init())
			// } else {
			// 	m.Next()
			// 	cmds = append(cmds, m.viewport.Init(), m.textarea.Focus())
			// }
		}
		switch m.state {
		case sidebarView:
			// apply update to models in view
			// cmds = append(cmds)
		}
	case responseMsg:
		m.responses++
		return m, waitForActivity(m.sub) // wait for next message
	case errMsg:
		m.err = msg
		return m, nil
	}
	return m, tea.Batch(cmds...)
}

func (m mainModel) View() string {
	var s string
	model := m.currentFocusedModel()
	// if m.state == sidebarView {
	s += lipgloss.JoinHorizontal(
		lipgloss.Center,
		docStyle.Render(m.sidebar.View()),
		fmt.Sprintf(
			"%s%s%s",
			m.viewport.View(),
			"\n",
			m.textarea.View(),
		),
	)
	// }
	s += helpStyle.Render(fmt.Sprintf("\ntab: focus next • n: new %s • q: exit\n", model))
	return s
}

func (m sidebarModel) View() string {
	return m.list.View()
}

func (m mainModel) currentFocusedModel() string {
	// if m.state == mainView {
	// 	return "main"
	// }
	return "main"
}

func (m mainModel) Next() {
	if m.index > 0 {
		m.index = 0
	} else {
		m.index++
	}
}

// chat notification
type responseMsg struct{}

func waitForActivity(sub chan struct{}) tea.Cmd {
	return func() tea.Msg {
		return responseMsg(<-sub)
	}
}

func init() {
	command.RootCmd.AddCommand(ChatCmd)
}
