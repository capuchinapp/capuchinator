package github

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// GitHub представляет работу с GitHub.
type GitHub struct {
	personalAccessToken string
	apiVersion          string
	packageURL          string
}

// New возвращает новый экземпляр GitHub.
func New(personalAccessToken string, apiVersion string, packageURL string) *GitHub {
	return &GitHub{
		personalAccessToken: personalAccessToken,
		apiVersion:          apiVersion,
		packageURL:          packageURL,
	}
}

func (g *GitHub) GetAvailableTags() ([]string, error) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, g.packageURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", "Bearer "+g.personalAccessToken)
	req.Header.Set("X-GitHub-Api-Version", g.apiVersion)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var packages []PackageInfo
	if err := json.Unmarshal(body, &packages); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	var tags []string
	for _, pkg := range packages {
		tags = append(tags, pkg.Metadata.Container.Tags...)
	}

	return tags, nil
}
