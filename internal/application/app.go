package application

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"

	"capuchinator/internal/domain"
	"capuchinator/internal/model"
	"capuchinator/internal/provider/docker"
	"capuchinator/internal/provider/github"
)

const (
	summaryWidth = 50

	compensationWidth  = 4
	compensationHeight = 2
)

type App struct {
	summary *model.Summary

	physicalWidth  int
	physicalHeight int

	leftStyle  lipgloss.Style
	rightStyle lipgloss.Style

	currentCmd tea.Model
	keys       domain.KeyMap
}

func New() (*App, error) {
	conf, err := LoadConfiguration()
	if err != nil {
		return nil, fmt.Errorf("load configuration: %v", err)
	}

	physicalWidth, physicalHeight, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return nil, fmt.Errorf("get terminal size: %v", err)
	}

	dockerService, err := docker.New(conf.DevMode)
	if err != nil {
		return nil, fmt.Errorf("new docker: %v", err)
	}

	theme := domain.NewTheme()

	summary := model.NewSummary(
		conf.DevMode,
		summaryWidth,
		theme,
		conf.Filename.ComposeBlue,
		conf.Filename.ComposeGreen,
		conf.Filename.NginxConf,
		conf.Filename.VictoriaMetrics,
		conf.Filename.Vector,
	)

	dic := NewDIC(
		conf.DevMode,
		summaryWidth,
		physicalWidth,
		physicalHeight,
		theme,
		summary,
		dockerService,
		github.New(conf.GitHub.PersonalAccessToken, conf.GitHub.APIVersion, conf.GitHub.PackagesURL),
	)

	return &App{
		summary: summary,

		physicalWidth:  physicalWidth,
		physicalHeight: physicalHeight,

		leftStyle: lipgloss.NewStyle().
			Width(summaryWidth).
			Height(physicalHeight-compensationHeight).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(theme.ColorGreen).
			Padding(0, 1),
		rightStyle: lipgloss.NewStyle().
			Width(physicalWidth-summaryWidth-compensationWidth).
			Height(physicalHeight-compensationHeight).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(theme.ColorGreen).
			Padding(0, 1),

		currentCmd: model.NewDir(dic),
		keys:       domain.NewKeyMap(),
	}, nil
}

func (a *App) Init() tea.Cmd {
	return a.currentCmd.Init()
}

func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	a.currentCmd, cmd = a.currentCmd.Update(msg)

	switch msg := msg.(type) {
	case model.NextCmdMsg:
		a.currentCmd = msg.NextCmd
		return a, a.currentCmd.Init()
	case tea.KeyMsg:
		if key.Matches(msg, a.keys.Quit) { //nolint:nestif // ignore
			if a.summary.GetDevMode() {
				output, err := exec.Command("docker", "rm", "-f", "capuchinator_app").CombinedOutput()
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println(string(output))

				output, err = exec.Command("docker", "rmi", "-f", "nginx:alpine").CombinedOutput()
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println(string(output))
			}

			return a, tea.Quit
		}
	}

	return a, cmd
}

func (a *App) View() string {
	return lipgloss.JoinHorizontal(
		lipgloss.Bottom,
		a.leftStyle.Render(a.summary.View()),
		a.rightStyle.Render(a.currentCmd.View()),
	)
}
