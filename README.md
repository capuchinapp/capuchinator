[![audit](https://github.com/capuchinapp/capuchinator/actions/workflows/audit.yml/badge.svg?branch=master)](https://github.com/capuchinapp/capuchinator/actions/workflows/audit.yml)

# Capuchinator

## Установка

```bash
cd /opt/capuchin
curl -L https://raw.githubusercontent.com/capuchinapp/capuchinator/refs/heads/master/scripts/install.sh
chmod +x install.sh
./install.sh v0.1.0
```

## Использование

1. Создайте символическую ссылку на файл конфига nginx

```bash
ln -sf /etc/nginx/sites-available/capuchin.ru /opt/capuchin/nginx_capuchin.conf
```

2. Настройка переменных окружения (используется инструмент [direnv](https://github.com/direnv/direnv))

```bash
echo 'export CAPUCHINATOR_GITHUB_PERSONAL_ACCESS_TOKEN="ghp_key"
export DOCKER_API_VERSION="1.43"' > .envrc
```

3. Запустите приложение

```bash
./capuchinator
```
