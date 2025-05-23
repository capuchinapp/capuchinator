# Capuchinator

## Установка

1. Создайте символическую ссылку на файл конфига nginx

```bash
ln -sf /etc/nginx/sites-available/capuchin.ru /opt/capuchin/nginx_capuchin.conf
```

2. Скопируйте файл `.envrc.example` с именем `/opt/capuchin/.envrc` и переопределите параметры

3. Запустите приложение
```bash
./capuchinator
```
