package model

import (
	"fmt"
	"os/exec"

	"capuchinator/internal/domain"
)

func NewExecStoppingCurrentDeploy(dic DIC) *Exec {
	summary := dic.GetSummary()

	return NewExec(dic, domain.ExecConfig{
		Name: "Stopping current deploy",

		StartFunc: func() domain.ExecResult {
			currentStrategy := summary.GetCurrentStrategy().String()
			currentVersion := summary.GetCurrentVersion()

			containerAPI := fmt.Sprintf(`capuchin_%s_%s_api`, currentStrategy, currentVersion)
			containerUI := fmt.Sprintf(`capuchin_%s_%s_ui`, currentStrategy, currentVersion)

			command := exec.Command(
				PathDocker,
				"stop",
				containerAPI,
				containerUI,
			)
			if dic.GetDevMode() {
				command = exec.Command(
					PathDocker,
					"ps",
				)
			}

			output, err := command.CombinedOutput()
			if err != nil {
				return domain.ExecResult{
					Status: domain.ExecResultStatusError,
					Err:    err,
					Output: string(output),
				}
			}

			return domain.ExecResult{
				Status: domain.ExecResultStatusSuccess,
				Output: "### Stopping current deploy:\n" + string(output),
			}
		},

		SuccessFunc: func() {
			summary.UpdateShutdownStopping(true)
		},

		ErrorFunc: func() {
			summary.UpdateShutdownStopping(false)
		},

		NextCmd: NewClearPGStat(dic),
	})
}
