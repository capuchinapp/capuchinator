package application

import (
	"capuchinator/internal/domain"
	"capuchinator/internal/model"
	"capuchinator/internal/provider/docker"
	"capuchinator/internal/provider/github"
)

type DIC struct {
	devMode bool

	summaryWidth int

	physicalWidth  int
	physicalHeight int

	theme *domain.Theme

	summary *model.Summary

	docker *docker.Docker
	gitHub *github.GitHub
}

func NewDIC(
	devMode bool,
	summaryWidth int,
	physicalWidth int,
	physicalHeight int,
	theme *domain.Theme,
	summary *model.Summary,
	dockerService *docker.Docker,
	gitHub *github.GitHub,
) *DIC {
	return &DIC{
		devMode: devMode,

		summaryWidth: summaryWidth,

		physicalWidth:  physicalWidth,
		physicalHeight: physicalHeight,

		theme: theme,

		summary: summary,

		docker: dockerService,
		gitHub: gitHub,
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

func (d *DIC) GetDocker() *docker.Docker {
	return d.docker
}

func (d *DIC) GetGitHub() *github.GitHub {
	return d.gitHub
}
