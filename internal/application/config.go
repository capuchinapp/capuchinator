package application

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

// Configuration представляет конфигурацию приложения.
type Configuration struct {
	DevMode  bool `env:"CAPUCHINATOR_DEV_MODE" envDefault:"false"`
	Filename Filename
	GitHub   GitHub
}

type Filename struct {
	// ComposeBlue содержит имя compose-файла для Blue.
	ComposeBlue string `env:"CAPUCHINATOR_FILENAME_COMPOSE_BLUE" envDefault:"compose.blue.yaml"`

	// ComposeGreen содержит имя compose-файла для Green.
	ComposeGreen string `env:"CAPUCHINATOR_FILENAME_COMPOSE_GREEN" envDefault:"compose.green.yaml"`

	// NginxConf содержит имя СИМВОЛИЧЕСКОЙ ССЫЛКИ на конфигурационный файл Nginx.
	NginxConf string `env:"CAPUCHINATOR_FILENAME_NGINX_CONF" envDefault:"nginx_capuchin.conf"`

	// VictoriaMetrics содержит имя файла конфига victoriametrics.
	VictoriaMetrics string `env:"CAPUCHINATOR_FILENAME_VICTORIAMETRICS" envDefault:"victoriametrics.yaml"`

	// Vector содержит имя файла конфига vector.
	Vector string `env:"CAPUCHINATOR_FILENAME_VECTOR" envDefault:"vector.yaml"`
}

// GitHub представляет конфигурацию GitHub.
type GitHub struct {
	// PersonalAccessToken содержит токен доступа к GitHub API.
	PersonalAccessToken string `env:"CAPUCHINATOR_GITHUB_PERSONAL_ACCESS_TOKEN,required"`

	// APIVersion содержит версию GitHub API.
	APIVersion string `env:"CAPUCHINATOR_GITHUB_API_VERSION" envDefault:"2022-11-28"`

	// PackagesURL содержит URL для получения списка доступных версий.
	PackagesURL string `env:"CAPUCHINATOR_GITHUB_PACKAGES_URL" envDefault:"https://api.github.com/users/capuchinapp/packages/container/cloud%2Fapi/versions"`
}

// LoadConfiguration возвращает новую конфигурацию приложения на основе переменных среды.
func LoadConfiguration() (*Configuration, error) {
	var config Configuration
	if err := env.Parse(&config); err != nil {
		return nil, fmt.Errorf("parse configuration: %v", err)
	}

	return &config, nil
}
