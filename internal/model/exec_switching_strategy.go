package model

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"capuchinator/internal/domain"
)

const (
	fileMode = 0644
)

func NewExecSwitchingStrategy(dic DIC) *Exec {
	summary := dic.GetSummary()

	return NewExec(dic, domain.ExecConfig{
		Name: "Switching strategy",

		StartFunc: func() domain.ExecResult {
			dir := summary.GetDir()
			devMode := dic.GetDevMode()
			currPortAPI, currPortUI := summary.GetCurrentPorts()
			nextPortAPI, nextPortUI := summary.GetNextPorts()

			pathNginx := path.Join(dir, summary.GetFilenameNginxConf())

			cmdTest := exec.Command(PathNginx, "-t")
			cmdReload := exec.Command(PathNginx, "-s", "reload")
			if devMode {
				cmdTest = exec.Command(PathDocker, "exec", TestContainerName, PathNginx, "-t")
				cmdReload = exec.Command(PathDocker, "exec", TestContainerName, PathNginx, "-s", "reload")
			}
			resNginx := switchNginx(pathNginx, currPortAPI, currPortUI, nextPortAPI, nextPortUI, cmdTest, cmdReload)
			if resNginx.Status == domain.ExecResultStatusError {
				return resNginx
			}

			return domain.ExecResult{
				Status: domain.ExecResultStatusSuccess,
				Output: fmt.Sprintf( //nolint:perfsprint // ignore
					"### Switching - Nginx:\n%s",
					resNginx.Output,
				),
			}
		},

		SuccessFunc: func() {
			summary.UpdateSwitchingNginx(true)
		},

		ErrorFunc: func() {
			summary.UpdateSwitchingNginx(false)
		},

		NextCmd: NewExecStoppingCurrentDeploy(dic),
	})
}

func switchNginx(
	filePath string,
	currPortAPI string,
	currPortUI string,
	nextPortAPI string,
	nextPortUI string,
	cmdTest *exec.Cmd,
	cmdReload *exec.Cmd,
) domain.ExecResult {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return domain.ExecResult{
			Status: domain.ExecResultStatusError,
			Err:    fmt.Errorf("switchNginx: failed to read file: %v", err),
		}
	}

	fileContent := string(content)

	const proxyPass = "proxy_pass http://127.0.0.1:" //#nosec G101 -- This is a false positive

	currAPIString := proxyPass + currPortAPI
	nextAPIString := proxyPass + nextPortAPI

	currUIString := proxyPass + currPortUI
	nextUIString := proxyPass + nextPortUI

	if !strings.Contains(fileContent, "#"+currAPIString) {
		fileContent = strings.Replace(fileContent, currAPIString, "#"+currAPIString, 1)
	}

	if strings.Contains(fileContent, "#"+nextAPIString) {
		fileContent = strings.Replace(fileContent, "#"+nextAPIString, nextAPIString, 1)
	}

	if !strings.Contains(fileContent, "#"+currUIString) {
		fileContent = strings.Replace(fileContent, currUIString, "#"+currUIString, 1)
	}

	if strings.Contains(fileContent, "#"+nextUIString) {
		fileContent = strings.Replace(fileContent, "#"+nextUIString, nextUIString, 1)
	}

	err = os.WriteFile(filePath, []byte(fileContent), fileMode) //#nosec G306 -- This is a false positive
	if err != nil {
		return domain.ExecResult{
			Status: domain.ExecResultStatusError,
			Err:    fmt.Errorf("switchNginx: failed to write file: %v", err),
		}
	}

	outputTest, err := cmdTest.CombinedOutput()
	if err != nil {
		return domain.ExecResult{
			Status: domain.ExecResultStatusError,
			Err:    err,
			Output: string(outputTest),
		}
	}

	outputReload, err := cmdReload.CombinedOutput()
	if err != nil {
		return domain.ExecResult{
			Status: domain.ExecResultStatusError,
			Err:    err,
			Output: string(outputReload),
		}
	}

	return domain.ExecResult{
		Status: domain.ExecResultStatusSuccess,
		Output: fmt.Sprintf(
			"### Test config:\n%s\n### Reload Nginx:\n%s",
			string(outputTest),
			string(outputReload),
		),
	}
}
