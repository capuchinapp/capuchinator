package model

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"regexp"
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
			nextStrategy := summary.GetNextStrategy()
			nextVersion := summary.GetNextVersion()
			currPortAPI, currPortUI := summary.GetCurrentPorts()
			nextPortAPI, nextPortUI := summary.GetNextPorts()

			pathVictoriaMetrics := path.Join(dir, summary.GetFilenameVictoriaMetrics())
			pathVector := path.Join(dir, summary.GetFilenameVector())
			pathNginx := path.Join(dir, summary.GetFilenameNginxConf())

			containerVM := "capuchin_ops_victoriametrics"
			if devMode {
				containerVM = TestContainerName
			}
			resVictoriaMetrics := switchVictoriaMetrics(pathVictoriaMetrics, nextStrategy, containerVM)
			if resVictoriaMetrics.Status == domain.ExecResultStatusError {
				return resVictoriaMetrics
			}

			containerVector := "capuchin_ops_vector"
			if devMode {
				containerVector = TestContainerName
			}
			resVector := switchVector(pathVector, nextStrategy, nextVersion, containerVector)
			if resVector.Status == domain.ExecResultStatusError {
				return resVector
			}

			cmdTest := exec.Command("nginx", "-t")
			cmdReload := exec.Command("nginx", "-s", "reload")
			if devMode {
				cmdTest = exec.Command("docker", "exec", TestContainerName, "nginx", "-t")
				cmdReload = exec.Command("docker", "exec", TestContainerName, "nginx", "-s", "reload")
			}
			resNginx := switchNginx(pathNginx, currPortAPI, currPortUI, nextPortAPI, nextPortUI, cmdTest, cmdReload)
			if resNginx.Status == domain.ExecResultStatusError {
				return resNginx
			}

			return domain.ExecResult{
				Status: domain.ExecResultStatusSuccess,
				Output: fmt.Sprintf(
					"### Switching - VictoriaMetrics:\n%s\n### Switching - Vector:\n%s\n### Switching - Nginx:\n%s",
					resVictoriaMetrics.Output,
					resVector.Output,
					resNginx.Output,
				),
			}
		},

		SuccessFunc: func() {
			summary.UpdateSwitchingVictoriaMetrics(true)
			summary.UpdateSwitchingVector(true)
			summary.UpdateSwitchingNginx(true)
		},

		ErrorFunc: func() {
			summary.UpdateSwitchingVictoriaMetrics(false)
			summary.UpdateSwitchingVector(false)
			summary.UpdateSwitchingNginx(false)
		},

		NextCmd: NewExecStoppingCurrentDeploy(dic),
	})
}

func switchVictoriaMetrics(filePath string, nextStrategy domain.Strategy, container string) domain.ExecResult {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return domain.ExecResult{
			Status: domain.ExecResultStatusError,
			Err:    fmt.Errorf("switchVictoriaMetrics: failed to read file: %v", err),
		}
	}

	fileContent := string(content)

	blueString := "- targets: ['capuchin_blue_api:8080'] # blue"
	greenString := "- targets: ['capuchin_green_api:8080'] # green"

	if nextStrategy == domain.StrategyBlue { //nolint:nestif // ignore
		if strings.Contains(fileContent, "#"+blueString) {
			fileContent = strings.Replace(fileContent, "#"+blueString, blueString, 1)
		}

		if !strings.Contains(fileContent, "#"+greenString) {
			fileContent = strings.Replace(fileContent, greenString, "#"+greenString, 1)
		}
	} else {
		if !strings.Contains(fileContent, "#"+blueString) {
			fileContent = strings.Replace(fileContent, blueString, "#"+blueString, 1)
		}

		if strings.Contains(fileContent, "#"+greenString) {
			fileContent = strings.Replace(fileContent, "#"+greenString, greenString, 1)
		}
	}

	err = os.WriteFile(filePath, []byte(fileContent), fileMode) //#nosec G306 -- This is a false positive
	if err != nil {
		return domain.ExecResult{
			Status: domain.ExecResultStatusError,
			Err:    fmt.Errorf("switchVictoriaMetrics: failed to write file: %v", err),
		}
	}

	output, err := exec.Command("docker", "restart", container).CombinedOutput()
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

func switchVector(filePath string, nextStrategy domain.Strategy, nextVersion string, container string) domain.ExecResult {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return domain.ExecResult{
			Status: domain.ExecResultStatusError,
			Err:    fmt.Errorf("switchVector: failed to read file: %v", err),
		}
	}

	substr := nextStrategy.String() + "_" + nextVersion
	re := regexp.MustCompile(`(capuchin_)[a-z]+_v\d+\.\d+\.\d+(_[a-z]+)`)
	fileContent := re.ReplaceAllString(string(content), "${1}"+substr+"${2}")

	err = os.WriteFile(filePath, []byte(fileContent), fileMode) //#nosec G306 -- This is a false positive
	if err != nil {
		return domain.ExecResult{
			Status: domain.ExecResultStatusError,
			Err:    fmt.Errorf("switchVector: failed to write file: %v", err),
		}
	}

	output, err := exec.Command("docker", "restart", container).CombinedOutput()
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

	currAPIString := "proxy_pass http://127.0.0.1:" + currPortAPI
	nextAPIString := "proxy_pass http://127.0.0.1:" + nextPortAPI

	currUIString := "proxy_pass http://127.0.0.1:" + currPortUI
	nextUIString := "proxy_pass http://127.0.0.1:" + nextPortUI

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
