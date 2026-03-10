package model

import (
	"errors"
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/docker/docker/api/types/container"

	"capuchinator/internal/domain"
	"capuchinator/internal/provider/docker"
)

const (
	maxRetries = 10
	retryDelay = 3 * time.Second
)

func NewExecHealthcheck(dic DIC) *Exec {
	summary := dic.GetSummary()
	dockerService := dic.GetDockerService()

	return NewExec(dic, domain.ExecConfig{
		Name: "Healthcheck",

		StartFunc: func() domain.ExecResult {
			nextVersion := summary.GetNextVersion()

			containers, err := dockerService.GetNextDeployContainers(nextVersion)
			if err != nil {
				return domain.ExecResult{
					Status: domain.ExecResultStatusError,
					Err:    err,
					Output: "failed to get next deploy containers",
				}
			}

			return healthcheck(dockerService, nextVersion, containers)
		},

		SuccessFunc: func() {
			summary.UpdateDeployHealthcheck(true)
		},

		ErrorFunc: func() {
			summary.UpdateDeployHealthcheck(false)
		},

		NextCmd: func() tea.Model {
			if summary.GetMode() == domain.ModeUpdate {
				return NewExecSwitchingStrategy(dic)
			}

			return NewComplete(dic.GetTheme())
		}(),
	})
}

func healthcheck(dockerService *docker.Docker, nextVersion string, containers []container.Summary) domain.ExecResult {
	var lastResult domain.ExecResult

	for attempt := 1; attempt <= maxRetries; attempt++ {
		var stateUI, stateAPI domain.ContainerState
		for _, ctr := range containers {
			state, err := dockerService.GetState(ctr.ID)
			if err != nil {
				return domain.ExecResult{
					Status: domain.ExecResultStatusError,
					Err:    err,
					Output: "failed to get container state",
				}
			}

			switch ctr.Image {
			case "ghcr.io/capuchinapp/cloud/ui:" + nextVersion:
				stateUI = state
			case "ghcr.io/capuchinapp/cloud/api:" + nextVersion:
				stateAPI = state
			}
		}

		if stateUI.Status == domain.ContainerStateStatusRunning &&
			stateUI.Health == domain.ContainerStateHealthHealthy &&
			stateAPI.Status == domain.ContainerStateStatusRunning &&
			stateAPI.Health == domain.ContainerStateHealthHealthy {
			return domain.ExecResult{
				Status: domain.ExecResultStatusSuccess,
				Err:    nil,
				Output: fmt.Sprintf(
					"###  UI State: %s (%s)\n### API State: %s (%s)",
					stateUI.Status,
					stateUI.Health,
					stateAPI.Status,
					stateAPI.Health,
				),
			}
		}

		lastResult = domain.ExecResult{
			Status: domain.ExecResultStatusError,
			Err:    errors.New("not ready"),
			Output: fmt.Sprintf(
				"###  UI State: %s (%s)\n### API State: %s (%s)\n(attempt %d/%d)",
				stateUI.Status,
				stateUI.Health,
				stateAPI.Status,
				stateAPI.Health,
				attempt,
				maxRetries,
			),
		}

		if attempt < maxRetries {
			time.Sleep(retryDelay)
		}
	}

	return lastResult
}
