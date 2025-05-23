package model

import (
	"capuchinator/internal/domain"
	"capuchinator/internal/provider/docker"
	"capuchinator/internal/provider/github"
)

type DIC interface {
	GetDevMode() bool

	GetSummaryWidth() int

	GetPhysicalWidth() int
	GetPhysicalHeight() int

	GetTheme() *domain.Theme

	GetSummary() *Summary

	GetDocker() *docker.Docker
	GetGitHub() *github.GitHub
}
