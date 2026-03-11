[![audit](https://github.com/capuchinapp/capuchinator/actions/workflows/audit.yml/badge.svg?branch=master)](https://github.com/capuchinapp/capuchinator/actions/workflows/audit.yml)

# Capuchinator

## Установка

```bash
curl -LsSf https://raw.githubusercontent.com/capuchinapp/capuchinator/refs/heads/master/scripts/install.sh | sh
```

## Использование

> Обратите внимание команды выполняются для каталога: `/opt/capuchin`

1. Создайте символическую ссылку на файл конфига nginx

```bash
ln -sf /etc/nginx/sites-available/capuchin.ru /opt/capuchin/nginx_capuchin.conf
```

2. Настройка переменных окружения (используется инструмент [direnv](https://github.com/direnv/direnv))

```bash
echo 'export CAPUCHINATOR_GITHUB_PERSONAL_ACCESS_TOKEN="ghp_key"
export DOCKER_API_VERSION="1.43"' > /opt/capuchin/.envrc
```

3. Запустите приложение

```bash
capuchinator
```
