package model

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"gopkg.in/yaml.v3"

	"capuchinator/internal/domain"
)

type Ports struct {
	dic     DIC
	summary *Summary
	theme   *domain.Theme
	spinner spinner.Model

	ports portsInfo
	err   error
}

type portsInfo struct {
	blue  ports
	green ports
}

type ports struct {
	api string
	ui  string
}

func NewPorts(dic DIC) *Ports {
	return &Ports{
		dic:     dic,
		summary: dic.GetSummary(),
		theme:   dic.GetTheme(),
		spinner: domain.NewSpinner(dic.GetTheme().StyleGreen),
	}
}

func (c *Ports) Init() tea.Cmd {
	return tea.Batch(c.processing, c.spinner.Tick)
}

func (c *Ports) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case StatusDone:
		if c.summary.GetMode() == domain.ModeInstall {
			c.summary.UpdateCurrentPorts("---", "---")
			c.summary.UpdateNextPorts(c.ports.blue.api, c.ports.blue.ui)

			return c, func() tea.Msg {
				return NextCmdMsg{
					NextCmd: NewAvailableTags(c.dic),
				}
			}
		}

		if c.summary.GetCurrentStrategy() == domain.StrategyBlue {
			c.summary.UpdateCurrentPorts(c.ports.blue.api, c.ports.blue.ui)
			c.summary.UpdateNextPorts(c.ports.green.api, c.ports.green.ui)

			return c, func() tea.Msg {
				return NextCmdMsg{
					NextCmd: NewAvailableTags(c.dic),
				}
			}
		}

		c.summary.UpdateCurrentPorts(c.ports.green.api, c.ports.green.ui)
		c.summary.UpdateNextPorts(c.ports.blue.api, c.ports.blue.ui)

		return c, func() tea.Msg {
			return NextCmdMsg{
				NextCmd: NewAvailableTags(c.dic),
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

func (c *Ports) View() string {
	prefix := "Searching ports from docker compose files"

	if c.err != nil {
		return prefix + "... " + c.theme.StyleRed.Render(
			fmt.Sprintf("something went wrong: %s", c.err),
		)
	}

	return fmt.Sprintf("%s %s", prefix, c.spinner.View())
}

func (c *Ports) processing() tea.Msg {
	dir := c.summary.GetDir()

	pathBlue := path.Join(dir, c.summary.GetFilenameComposeBlue())
	pathGreen := path.Join(dir, c.summary.GetFilenameComposeGreen())

	var err error

	c.ports.blue.api, c.ports.blue.ui, err = findPorts(domain.StrategyBlue, pathBlue)
	if err != nil {
		return StatusError{
			fmt.Errorf("failed to find blue ports: %v", err),
		}
	}

	c.ports.green.api, c.ports.green.ui, err = findPorts(domain.StrategyGreen, pathGreen)
	if err != nil {
		return StatusError{
			fmt.Errorf("failed to find green ports: %v", err),
		}
	}

	return StatusDone{true}
}

func findPorts(strategy domain.Strategy, pathFile string) (api string, ui string, err error) {
	data, err := os.ReadFile(pathFile)
	if err != nil {
		return "", "", fmt.Errorf("failed to read file: %v", err)
	}

	if strategy == domain.StrategyBlue {
		cfg := ConfigBlue{}

		err = yaml.Unmarshal(data, &cfg)
		if err != nil {
			return "", "", fmt.Errorf("yaml unmarshal: %v", err)
		}

		api = strings.Split(cfg.Services.API.Ports[0], ":")[0]
		ui = strings.Split(cfg.Services.UI.Ports[0], ":")[0]

		return api, ui, nil
	}

	cfg := ConfigGreen{}

	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return "", "", fmt.Errorf("yaml unmarshal: %v", err)
	}

	api = strings.Split(cfg.Services.API.Ports[0], ":")[0]
	ui = strings.Split(cfg.Services.UI.Ports[0], ":")[0]

	return api, ui, nil
}

type ConfigBlue struct {
	Services ServicesBlue `yaml:"services"`
}

type ConfigGreen struct {
	Services ServicesGreen `yaml:"services"`
}

type ServicesBlue struct {
	API Service `yaml:"capuchin_blue_api"`
	UI  Service `yaml:"capuchin_blue_ui"`
}

type ServicesGreen struct {
	API Service `yaml:"capuchin_green_api"`
	UI  Service `yaml:"capuchin_green_ui"`
}

type Service struct {
	Ports []string `yaml:"ports"`
}
