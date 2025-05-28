package application

import (
	"capuchinator/internal/domain"
	"capuchinator/internal/model"
	"capuchinator/internal/provider/docker"
	"capuchinator/internal/provider/github"
)

type DICConfig struct {
	DevMode bool

	SummaryWidth int

	PhysicalWidth  int
	PhysicalHeight int

	Theme *domain.Theme

	Summary *model.Summary

	DockerService *docker.Docker
	GitHubService *github.GitHub
}

type DIC struct {
	devMode bool

	summaryWidth int

	physicalWidth  int
	physicalHeight int

	theme *domain.Theme

	summary *model.Summary

	dockerService *docker.Docker
	gitHubService *github.GitHub
}

func NewDIC(cfg DICConfig) *DIC {
	return &DIC{
		devMode: cfg.DevMode,

		summaryWidth: cfg.SummaryWidth,

		physicalWidth:  cfg.PhysicalWidth,
		physicalHeight: cfg.PhysicalHeight,

		theme: cfg.Theme,

		summary: cfg.Summary,

		dockerService: cfg.DockerService,
		gitHubService: cfg.GitHubService,
	}
}

func (d *DIC) GetDevMode() bool {
	return d.devMode
}

func (d *DIC) GetSummaryWidth() int {
	return d.summaryWidth
}

func (d *DIC) GetPhysicalWidth() int {
	return d.physicalWidth
}

func (d *DIC) GetPhysicalHeight() int {
	return d.physicalHeight
}

func (d *DIC) GetTheme() *domain.Theme {
	return d.theme
}

func (d *DIC) GetSummary() *model.Summary {
	return d.summary
}

func (d *DIC) GetDockerService() *docker.Docker {
	return d.dockerService
}

func (d *DIC) GetGitHubService() *github.GitHub {
	return d.gitHubService
}
