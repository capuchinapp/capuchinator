package model

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"capuchinator/internal/domain"
)

func NewExecCheckingRequests(dic DIC) *Exec {
	summary := dic.GetSummary()

	return NewExec(dic, domain.ExecConfig{
		Name: "Checking requests",

		StartFunc: func() domain.ExecResult {
			portAPI, portUI := summary.GetNextPorts()
			if dic.GetDevMode() {
				portAPI = "8585"
				portUI = "8585"
			}

			resAPI := checkRequest(portAPI)
			if resAPI.Status == domain.ExecResultStatusError {
				return resAPI
			}

			resUI := checkRequest(portUI)
			if resUI.Status == domain.ExecResultStatusError {
				return resUI
			}

			return domain.ExecResult{
				Status: domain.ExecResultStatusSuccess,
				Output: fmt.Sprintf(
					"### Request to api (http://localhost:%s):\n%s\n### Request to ui (http://localhost:%s):\n%s",
					portAPI,
					resAPI.Output,
					portUI,
					resUI.Output,
				),
			}
		},

		SuccessFunc: func() {
			summary.UpdateDeployCheckingRequests(true)
		},

		ErrorFunc: func() {
			summary.UpdateDeployCheckingRequests(false)
		},

		NextCmd: func() tea.Model {
			if summary.GetMode() == domain.ModeUpdate {
				return NewExecSwitchingStrategy(dic)
			}

			return NewComplete(dic.GetTheme())
		}(),
	})
}

func checkRequest(port string) domain.ExecResult {
	cmd := exec.Command(PathCurl, "-I", "http://localhost:"+port) //#nosec G204 -- This is a false positive

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	output, err := cmd.Output()
	if err != nil {
		out := stderr.String()
		if strings.Contains(out, "curl:") {
			lines := strings.Split(out, "curl:")
			out = "curl:" + lines[1]
		}

		return domain.ExecResult{
			Status: domain.ExecResultStatusError,
			Err:    err,
			Output: out,
		}
	}

	return domain.ExecResult{
		Status: domain.ExecResultStatusSuccess,
		Output: string(output),
	}
}
