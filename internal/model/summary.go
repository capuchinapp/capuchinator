package model

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"capuchinator/internal/domain"
)

type SummaryConfig struct {
	AppVersion string
	DevMode    bool

	Width int

	Theme *domain.Theme

	FilenameComposeBlue     string
	FilenameComposeGreen    string
	FilenameNginxConf       string
	FilenameVictoriaMetrics string
	FilenameVector          string
}

type Summary struct {
	appVersion string
	devMode    bool
	theme      *domain.Theme

	dirMaxWidth int

	mode domain.Mode

	requirementsCurlVersion          string
	requirementsDockerVersion        string
	requirementsDockerComposeVersion string
	requirementsNginxVersion         string

	filenameComposeBlue     string
	filenameComposeGreen    string
	filenameNginxConf       string
	filenameVictoriaMetrics string
	filenameVector          string

	currentDir      string
	currentVersion  string
	currentStrategy domain.Strategy
	currentPortAPI  string
	currentPortUI   string

	nextVersion  string
	nextStrategy domain.Strategy
	nextPortAPI  string
	nextPortUI   string

	deployLaunching        *bool
	deployCheckingLogs     *bool
	deployCheckingRequests *bool

	switchingVictoriaMetrics *bool
	switchingVector          *bool
	switchingNginx           *bool

	shutdownStopping *bool

	styles styles
}

type styles struct {
	category lipgloss.Style
	title    lipgloss.Style
	text     lipgloss.Style
}

func NewSummary(cfg SummaryConfig) *Summary {
	marginCompensation := 3

	return &Summary{
		appVersion: cfg.AppVersion,
		devMode:    cfg.DevMode,
		theme:      cfg.Theme,

		dirMaxWidth: cfg.Width - marginCompensation,

		mode: "???",

		requirementsCurlVersion:          "???",
		requirementsDockerVersion:        "???",
		requirementsDockerComposeVersion: "???",
		requirementsNginxVersion:         "???",

		filenameComposeBlue:     cfg.FilenameComposeBlue,
		filenameComposeGreen:    cfg.FilenameComposeGreen,
		filenameNginxConf:       cfg.FilenameNginxConf,
		filenameVictoriaMetrics: cfg.FilenameVictoriaMetrics,
		filenameVector:          cfg.FilenameVector,

		currentDir:      "???",
		currentVersion:  "???",
		currentStrategy: "???",
		currentPortAPI:  "???",
		currentPortUI:   "???",

		nextVersion:  "???",
		nextStrategy: "???",
		nextPortAPI:  "???",
		nextPortUI:   "???",

		styles: styles{
			category: lipgloss.NewStyle().
				Foreground(cfg.Theme.ColorWhite).
				Bold(true).
				Transform(strings.ToUpper).
				MarginTop(1),

			title: lipgloss.NewStyle().
				Foreground(cfg.Theme.ColorGreen).
				Bold(true),

			text: lipgloss.NewStyle().
				Foreground(cfg.Theme.ColorYellow).
				Bold(true),
		},
	}
}

func (s *Summary) View() string {
	title := lipgloss.NewStyle().
		Foreground(s.theme.ColorOrange).
		Bold(true).
		Transform(strings.ToUpper).
		Render("🛠️ Capuchinator")
	version := lipgloss.NewStyle().
		Foreground(s.theme.ColorWhite).
		Render(s.appVersion)
	header := lipgloss.NewStyle().
		MarginBottom(1).
		Render(title + " v" + version)

	devModeStr := s.styles.text.Render("off")
	if s.devMode {
		devModeStr = lipgloss.NewStyle().
			Foreground(s.theme.ColorRed).
			Bold(true).
			Render("on")
	}

	var deploy string
	if s.mode == domain.ModeInstall || s.mode == domain.ModeUpdate {
		deploy = fmt.Sprintf(
			"%s\n%s\n%s\n%s",
			s.styles.category.Render("Deploy ("+s.nextVersion+")"),
			s.styles.title.Render("Launching:           ")+s.boolToIcon(s.deployLaunching),
			s.styles.title.Render("Checking - Logs:     ")+s.boolToIcon(s.deployCheckingLogs),
			s.styles.title.Render("Checking - Requests: ")+s.boolToIcon(s.deployCheckingRequests),
		)
	}

	var switchStrategy string
	var shutdown string
	if s.mode == domain.ModeUpdate {
		switchStrategy = fmt.Sprintf(
			"%s\n%s\n%s\n%s",
			s.styles.category.Render("Switch strategy"),
			s.styles.title.Render("Switching - VictoriaMetrics: ")+s.boolToIcon(s.switchingVictoriaMetrics),
			s.styles.title.Render("Switching - Vector:          ")+s.boolToIcon(s.switchingVector),
			s.styles.title.Render("Switching - Nginx:           ")+s.boolToIcon(s.switchingNginx),
		)

		shutdown = fmt.Sprintf(
			"%s\n%s",
			s.styles.category.Render("Shutdown ("+s.currentVersion+")"),
			s.styles.title.Render("Stopping the old version: ")+s.boolToIcon(s.shutdownStopping),
		)
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		s.styles.title.Render("DevMode:   ")+devModeStr,
		s.styles.title.Render("Mode:      ")+s.styles.text.Render(s.mode.String()),

		s.styles.category.Render("Requirements"),
		s.styles.title.Render("curl:           ")+s.styles.text.Render(s.requirementsCurlVersion),
		s.styles.title.Render("docker:         ")+s.styles.text.Render(s.requirementsDockerVersion),
		s.styles.title.Render("docker compose: ")+s.styles.text.Render(s.requirementsDockerComposeVersion),
		s.styles.title.Render("nginx:          ")+s.styles.text.Render(s.requirementsNginxVersion),

		s.styles.category.Render("Directory"),
		s.styles.text.Render(splitString(s.currentDir, s.dirMaxWidth)),

		s.styles.category.Render("Files"),
		s.styles.title.Render("compose.blue.yaml:    ")+s.styles.text.Render(s.filenameComposeBlue),
		s.styles.title.Render("compose.green.yaml:   ")+s.styles.text.Render(s.filenameComposeGreen),
		s.styles.title.Render("nginx.conf (symlink): ")+s.styles.text.Render(s.filenameNginxConf),
		s.styles.title.Render("victoriametrics.yaml: ")+s.styles.text.Render(s.filenameVictoriaMetrics),
		s.styles.title.Render("vector.yaml:          ")+s.styles.text.Render(s.filenameVector),

		s.styles.category.Render("Deploy strategy"),
		s.styles.title.Render("Version:    ")+s.styles.text.Render(s.currentVersion)+" >> "+s.styles.text.Render(s.nextVersion),
		s.styles.title.Render("Strategy:   ")+s.styles.text.Render(s.currentStrategy.String())+" >> "+s.styles.text.Render(s.nextStrategy.String()),
		s.styles.title.Render("Port - API: ")+s.styles.text.Render(s.currentPortAPI)+" >> "+s.styles.text.Render(s.nextPortAPI),
		s.styles.title.Render("Port - UI:  ")+s.styles.text.Render(s.currentPortUI)+" >> "+s.styles.text.Render(s.nextPortUI),

		deploy,
		switchStrategy,
		shutdown,
	)
}

func (s *Summary) GetDevMode() bool {
	return s.devMode
}

func (s *Summary) GetMode() domain.Mode {
	return s.mode
}

func (s *Summary) GetDir() string {
	return s.currentDir
}

func (s *Summary) GetFilenameComposeBlue() string {
	return s.filenameComposeBlue
}

func (s *Summary) GetFilenameComposeGreen() string {
	return s.filenameComposeGreen
}

func (s *Summary) GetFilenameNginxConf() string {
	return s.filenameNginxConf
}

func (s *Summary) GetFilenameVictoriaMetrics() string {
	return s.filenameVictoriaMetrics
}

func (s *Summary) GetFilenameVector() string {
	return s.filenameVector
}

func (s *Summary) GetCurrentVersion() string {
	return s.currentVersion
}

func (s *Summary) GetCurrentStrategy() domain.Strategy {
	return s.currentStrategy
}

func (s *Summary) GetCurrentPorts() (api string, ui string) {
	return s.currentPortAPI, s.currentPortUI
}

func (s *Summary) GetNextVersion() string {
	return s.nextVersion
}

func (s *Summary) GetNextStrategy() domain.Strategy {
	return s.nextStrategy
}

func (s *Summary) GetNextPorts() (api string, ui string) {
	return s.nextPortAPI, s.nextPortUI
}

func (s *Summary) UpdateDir(value string) {
	s.currentDir = value
}

func (s *Summary) UpdateMode(value domain.Mode) {
	s.mode = value
}

func (s *Summary) UpdateRequirementsCurlVersion(value string) {
	s.requirementsCurlVersion = value
}

func (s *Summary) UpdateRequirementsDockerVersion(value string) {
	s.requirementsDockerVersion = value
}

func (s *Summary) UpdateRequirementsDockerComposeVersion(value string) {
	s.requirementsDockerComposeVersion = value
}

func (s *Summary) UpdateRequirementsNginxVersion(value string) {
	s.requirementsNginxVersion = value
}

func (s *Summary) UpdateCurrentVersion(value string) {
	s.currentVersion = value
}

func (s *Summary) UpdateCurrentStrategy(value domain.Strategy) {
	s.currentStrategy = value
}

func (s *Summary) UpdateCurrentPorts(api string, ui string) {
	s.currentPortAPI = api
	s.currentPortUI = ui
}

func (s *Summary) UpdateNextVersion(value string) {
	s.nextVersion = value
}

func (s *Summary) UpdateNextStrategy(value domain.Strategy) {
	s.nextStrategy = value
}

func (s *Summary) UpdateNextPorts(api string, ui string) {
	s.nextPortAPI = api
	s.nextPortUI = ui
}

func (s *Summary) UpdateDeployLaunching(value bool) {
	s.deployLaunching = &value
}

func (s *Summary) UpdateDeployCheckingLogs(value bool) {
	s.deployCheckingLogs = &value
}

func (s *Summary) UpdateDeployCheckingRequests(value bool) {
	s.deployCheckingRequests = &value
}

func (s *Summary) UpdateSwitchingVictoriaMetrics(value bool) {
	s.switchingVictoriaMetrics = &value
}

func (s *Summary) UpdateSwitchingVector(value bool) {
	s.switchingVector = &value
}

func (s *Summary) UpdateSwitchingNginx(value bool) {
	s.switchingNginx = &value
}

func (s *Summary) UpdateShutdownStopping(value bool) {
	s.shutdownStopping = &value
}

func (s *Summary) boolToIcon(b *bool) string {
	if b == nil {
		return s.styles.text.Render("???")
	}

	if *b {
		return s.theme.StyleGreen.Render("✅")
	}

	return s.theme.StyleRed.Render("❌")
}

func splitString(s string, width int) string {
	var result []string

	for len(s) > width {
		result = append(result, s[:width])
		s = s[width:]
	}

	if len(s) > 0 {
		result = append(result, s)
	}

	return strings.Join(result, "\n")
}
