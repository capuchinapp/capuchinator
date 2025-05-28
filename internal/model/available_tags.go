package model

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"

	"capuchinator/internal/domain"
)

type AvailableTags struct {
	dic     DIC
	spinner spinner.Model

	tags []string
	err  error
}

func NewAvailableTags(dic DIC) *AvailableTags {
	return &AvailableTags{
		dic:     dic,
		spinner: domain.NewSpinner(dic.GetTheme().StyleGreen),
	}
}

func (c *AvailableTags) Init() tea.Cmd {
	return tea.Batch(c.getAvailableTags, c.spinner.Tick)
}

func (c *AvailableTags) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case StatusDone:
		return c, func() tea.Msg {
			return NextCmdMsg{
				NextCmd: NewNextDeploy(c.dic, c.tags),
			}
		}
	case StatusError:
		c.err = msg
		return c, nil
	default:
		var cmd tea.Cmd
		c.spinner, cmd = c.spinner.Update(msg)
		return c, cmd
	}
}

func (c *AvailableTags) View() string {
	prefix := "Getting available tags"

	if c.err != nil {
		return prefix + "... " + c.dic.GetTheme().StyleRed.Render(
			fmt.Sprintf("something went wrong: %s", c.err),
		)
	}

	return fmt.Sprintf("%s %s", prefix, c.spinner.View())
}

func (c *AvailableTags) getAvailableTags() tea.Msg {
	tags, err := c.dic.GetGitHubService().GetAvailableTags()
	if err != nil {
		return StatusError{
			fmt.Errorf("failed to get available tags: %v", err),
		}
	}

	c.tags = tags

	return StatusDone{true}
}
