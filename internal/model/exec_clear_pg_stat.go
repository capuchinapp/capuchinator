package model

import (
	"os/exec"

	"capuchinator/internal/domain"
)

func NewExecClearPGStat(dic DIC) *Exec {
	return NewExec(dic, domain.ExecConfig{
		Name: "Clear pg_stat_statements",

		StartFunc: func() domain.ExecResult {
			command := exec.Command(
				PathDocker,
				"exec",
				"-it",
				"capuchin_postgres",
				"psql",
				"-U",
				"postgres",
				"-d",
				"postgres",
				"-c",
				`"SELECT pg_stat_statements_reset();"`,
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
				Output: "### Clear pg_stat_statements:\n" + string(output),
			}
		},

		SuccessFunc: func() {
			// No operation
		},

		ErrorFunc: func() {
			// No operation
		},

		NextCmd: NewComplete(dic.GetTheme()),
	})
}
