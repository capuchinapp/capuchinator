[![audit](https://github.com/capuchinapp/capuchinator/actions/workflows/audit.yml/badge.svg?branch=master)](https://github.com/capuchinapp/capuchinator/actions/workflows/audit.yml)

# Capuchinator

## Установка

```bash
curl -o /opt/capuchin/capuchinator_linux_amd64.tar.gz https://github.com/capuchinapp/capuchinator/releases/download/v0.1.0/capuchinator_0.1.0_linux_amd64.tar.gz?token=TOKEN
tar -xzf *.tar.gz
chmod +x ./capuchinator
```

## Использование

1. Создайте символическую ссылку на файл конфига nginx

```bash
ln -sf /etc/nginx/sites-available/capuchin.ru /opt/capuchin/nginx_capuchin.conf
```

2. Настройка переменных окружения

```bash
echo 'export CAPUCHINATOR_GITHUB_PERSONAL_ACCESS_TOKEN="ghp_key"
export DOCKER_API_VERSION="1.43"' > .envrc
```

3. Запустите приложение

```bash
./capuchinator
```
