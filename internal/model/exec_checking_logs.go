package model

import (
	"fmt"
	"os/exec"

	"capuchinator/internal/domain"
)

func NewExecCheckingLogs(dic DIC) *Exec {
	summary := dic.GetSummary()

	return NewExec(dic, domain.ExecConfig{
		Name: "Checking logs",

		StartFunc: func() domain.ExecResult {
			devMode := dic.GetDevMode()
			strategy := summary.GetNextStrategy()
			version := summary.GetNextVersion()

			containerAPI := fmt.Sprintf(`capuchin_%s_%s_%s`, strategy, version, "api")
			if devMode {
				containerAPI = TestContainerName
			}
			logsAPI := getLogs(containerAPI)
			if logsAPI.Status == domain.ExecResultStatusError {
				return logsAPI
			}

			containerUI := fmt.Sprintf(`capuchin_%s_%s_%s`, strategy, version, "ui")
			if devMode {
				containerUI = TestContainerName
			}
			logsUI := getLogs(containerUI)
			if logsUI.Status == domain.ExecResultStatusError {
				return logsUI
			}

			return domain.ExecResult{
				Status: domain.ExecResultStatusSuccess,
				Output: fmt.Sprintf(
					"### Logs api:\n%s\n### Logs ui:\n%s",
					logsAPI.Output,
					logsUI.Output,
				),
			}
		},

		SuccessFunc: func() {
			summary.UpdateDeployCheckingLogs(true)
		},

		ErrorFunc: func() {
			summary.UpdateDeployCheckingLogs(false)
		},

		NextCmd: NewExecCheckingRequests(dic),
	})
}

func getLogs(container string) domain.ExecResult {
	output, err := exec.Command(PathDocker, "logs", container).CombinedOutput()
	if err != nil {
		return domain.ExecResult{
			Status: domain.ExecResultStatusError,
			Err:    err,
			Output: string(output),
		}
	}

	return domain.ExecResult{
		Status: domain.ExecResultStatusSuccess,
		Output: string(output),
	}
}
