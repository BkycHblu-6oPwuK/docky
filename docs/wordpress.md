# Конфигурация под битрикс

nginx + php (5.6, 7.1, 7.4, 8.1, 8.2, 8.3, 8.4, 8.5) + mysql + node 24 версии
а так же redis|memcached - при самостоятельной публикации

## Шаги публикации docker-compose.yml

Шаги:

0. Выбор фреймворка - bitrix
1. Выбрать версию php
2. Выбрать версию mysql - 5.7 или 8.0
3. Устанавливать ли node.js - Y или N, если установите Y, то будет создан сервис с node js 24 версии
    1. Если будет устанавливаться node.js, то нужно указать корневую директорию для него, то есть директория содержащая файл package.json. Путь указывается относительно корня сайта - local/js/vite или пустое поле если package.json в корне сайта.

После этого в директории где выполнялась команда появится docker-compose.yml файл с настроенными сервисами.

## Cron

Если необходимо добавить задания, то сделайте публикацию файлов с заданиями

```bash
docky publish --file cron_tasks
```

Запись заданий осуществляйте в:
- `_conf/app/cron/docky` - для пользователя сайта
- `_conf/app/cron/root` - для root пользователя

Полностью отключить cron в контейнере можно установив переменную окружения ```CRON_DISABLED=true|1``` в docker-compose.yml и перезапустить контейнеры.

```yaml
services:
  app:
    environment:
      - CRON_DISABLED=true|1
```


В редких случаях, если вы уже пользовались docky до этого и у вас не отключился cron, то выполните команду:

```bash
docky build
```

## Почта

Для тестирования почты используйте сервис mailhog.

Если сервиса нет в docker-compose.yml, то опубликуйте его:

```bash
docky publish --service mailhog
```

И настройте smtp отправку почты на этот сервис:

- Сервер - mailhog
- Порт - 1025
- Без авторизации

Проверяйте письма на странице - http://localhost:8025

----

Для отправки почты через  ```msmtp```

Вам необходимо настроить файл msmtprc в ``` _docker/app/msmtprc ```, или создайте свой файл и пробросьте его с помощью volumes в docker-compose.yml:

```yaml
- ${CONF_PATH}/app/msmtprc:/home/docky/.msmtprc
```

Привер такого файла конфигурации:

```
defaults
tls on
auth on
keepbcc on
tls_certcheck off
logfile /home/docky/msmtp.log

account yandex
port 587
host smtp.yandex.ru
user <your-email>@yandex.ru
password <your-app-password>
from <your-email>@yandex.ru
tls_starttls on

account default : yandex
```

Если почта не отправляется или в проверке системы написано что почта не работает, то проверьте логи ```msmtp``` в контейнере, которые находятся в файле ```/home/docky/msmtp.log```. Вероятнее всего произошла ошибка авторизации или почтовый сервис отклюнил отправку из-за подозрений в спаме.
