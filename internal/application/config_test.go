package application

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfiguration(t *testing.T) {
	tests := []struct {
		name       string
		setEnvFunc func()
		want       *Configuration
		wantErr    error
	}{
		{
			name: "Should load configuration from env with default values",
			setEnvFunc: func() {
				t.Setenv("CAPUCHINATOR_GITHUB_PERSONAL_ACCESS_TOKEN", "ghp_key")
			},
			want: &Configuration{
				DevMode: func() bool {
					// Это надо из-за .envrc и запуска тестов через make test

					fromEnv := os.Getenv("CAPUCHINATOR_DEV_MODE")
					if fromEnv == "" {
						return false
					}

					return fromEnv == "true"
				}(),
				Filename: Filename{
					ComposeBlue:     "compose.blue.yaml",
					ComposeGreen:    "compose.green.yaml",
					NginxConf:       "nginx_capuchin.conf",
					VictoriaMetrics: "victoriametrics.yaml",
					Vector:          "vector.yaml",
				},
				GitHub: GitHub{
					PersonalAccessToken: "ghp_key",
					APIVersion:          "2022-11-28",
					PackagesURL:         "https://api.github.com/users/capuchinapp/packages/container/cloud%2Fapi/versions",
				},
			},
			wantErr: nil,
		},
		{
			name: "Should load configuration from env",
			setEnvFunc: func() {
				t.Setenv("CAPUCHINATOR_DEV_MODE", "true")

				t.Setenv("CAPUCHINATOR_FILENAME_COMPOSE_BLUE", "compose-blue")
				t.Setenv("CAPUCHINATOR_FILENAME_COMPOSE_GREEN", "compose-green")
				t.Setenv("CAPUCHINATOR_FILENAME_NGINX_CONF", "nginx-conf")
				t.Setenv("CAPUCHINATOR_FILENAME_VICTORIAMETRICS", "victoriametrics")
				t.Setenv("CAPUCHINATOR_FILENAME_VECTOR", "vector")

				t.Setenv("CAPUCHINATOR_GITHUB_PERSONAL_ACCESS_TOKEN", "ghp_key")
				t.Setenv("CAPUCHINATOR_GITHUB_API_VERSION", "2022-11-29")
				t.Setenv("CAPUCHINATOR_GITHUB_PACKAGES_URL", "https://versions")
			},
			want: &Configuration{
				DevMode: true,
				Filename: Filename{
					ComposeBlue:     "compose-blue",
					ComposeGreen:    "compose-green",
					NginxConf:       "nginx-conf",
					VictoriaMetrics: "victoriametrics",
					Vector:          "vector",
				},
				GitHub: GitHub{
					PersonalAccessToken: "ghp_key",
					APIVersion:          "2022-11-29",
					PackagesURL:         "https://versions",
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setEnvFunc()
			got, err := LoadConfiguration()
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
