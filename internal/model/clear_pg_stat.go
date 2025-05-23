package model

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"

	"capuchinator/internal/domain"
)

type ClearPGStat struct {
	dic DIC

	confirm string

	lg   *lipgloss.Renderer
	form *huh.Form
}

func NewClearPGStat(dic DIC) *ClearPGStat {
	c := &ClearPGStat{
		dic: dic,
		lg:  lipgloss.DefaultRenderer(),
	}

	c.form = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Clear PG Stat").
				Options(huh.NewOptions(domain.ConfirmNo.String(), domain.ConfirmYes.String())...).
				Value(&c.confirm),
		).WithShowHelp(false),
	)

	return c
}

func (c *ClearPGStat) Init() tea.Cmd {
	return c.form.Init()
}

func (c *ClearPGStat) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	form, cmd := c.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		c.form = f
		cmds = append(cmds, cmd)
	}

	if c.form.State == huh.StateCompleted {
		if domain.Confirm(c.confirm) == domain.ConfirmYes {
			cmds = append(cmds, func() tea.Msg {
				return NextCmdMsg{
					NextCmd: NewExecClearPGStat(c.dic),
				}
			})

			return c, tea.Batch(cmds...)
		}

		cmds = append(cmds, func() tea.Msg {
			return NextCmdMsg{
				NextCmd: NewComplete(c.dic.GetTheme()),
			}
		})
	}

	return c, tea.Batch(cmds...)
}

func (c *ClearPGStat) View() string {
	v := strings.TrimSuffix(c.form.View(), "\n\n")

	return c.lg.NewStyle().Render(v)
}
