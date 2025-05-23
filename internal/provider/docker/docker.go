package docker

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/hashicorp/go-version"

	"capuchinator/internal/domain"
)

// Docker представляет работу с Docker.
type Docker struct {
	devMode bool

	cli *client.Client
}

// New возвращает новый экземпляр Docker.
func New(devMode bool) (*Docker, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, fmt.Errorf("client.NewClientWithOpts: %v", err)
	}

	return &Docker{
		devMode: devMode,

		cli: cli,
	}, nil
}

func (d *Docker) GetCurrentDeploy() (currVersion string, currStrategy domain.Strategy, err error) {
	containers, err := d.cli.ContainerList(context.Background(), container.ListOptions{})
	if err != nil {
		return "", "", fmt.Errorf("failed to list containers: %v", err)
	}

	if d.devMode {
		containers = []container.Summary{
			{
				Image: "ghcr.io/capuchinapp/cloud/ui:v0.8.0",
				Names: []string{"/capuchin_blue_v0.8.0_ui"},
			},
			{
				Image: "ghcr.io/capuchinapp/cloud/api:v0.8.0",
				Names: []string{"/capuchin_blue_v0.8.0_api"},
			},
		}
	}

	currVersion, err = getCurrVersion(containers)
	if err != nil {
		return "", "", fmt.Errorf("failed to get current version: %v", err)
	}

	currStrategy, err = getCurrStrategy(containers)
	if err != nil {
		return "", "", fmt.Errorf("failed to get current strategy: %v", err)
	}

	return currVersion, currStrategy, nil
}

func getCurrVersion(containers []container.Summary) (string, error) {
	versions := make(map[string]struct{})
	for _, ctr := range containers {
		if strings.HasPrefix(ctr.Image, "ghcr.io/capuchinapp/cloud") {
			v := strings.Split(ctr.Image, ":")[1]
			versions[v] = struct{}{}
		}
	}
	if len(versions) == 0 {
		return "", errors.New("no ghcr.io/capuchinapp/cloud containers found")
	}
	if len(versions) > 1 {
		return "", errors.New("multiple ghcr.io/capuchinapp/cloud containers found")
	}

	var currVersion string
	for k := range versions {
		currVersion = k
		break
	}

	v, err := version.NewSemver(currVersion)
	if err != nil {
		return "", fmt.Errorf("failed to parse version: %v", err)
	}

	return "v" + v.String(), nil
}

func getCurrStrategy(containers []container.Summary) (domain.Strategy, error) {
	strategies := make(map[domain.Strategy]struct{})
	for _, ctr := range containers {
		for _, name := range ctr.Names {
			if strings.HasPrefix(name, "/capuchin_blue") {
				strategies[domain.StrategyBlue] = struct{}{}
				continue
			}

			if strings.HasPrefix(name, "/capuchin_green") {
				strategies[domain.StrategyGreen] = struct{}{}
			}
		}
	}
	if len(strategies) == 0 {
		return "", errors.New("no strategy found")
	}
	if len(strategies) > 1 {
		return "", errors.New("multiple strategies found")
	}

	var currStrategy domain.Strategy
	for k := range strategies {
		currStrategy = k
		break
	}

	return currStrategy, nil
}
