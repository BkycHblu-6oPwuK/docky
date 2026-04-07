# Конфигурация под yii3

nginx + php (8.2, 8.3, 8.4, 8.5) + mysql|mariadb|postgres|sqlite
а так же redis|memcached

## Шаги публикации docker-compose.yml

Шаги:

0. Выбор фреймворка - yii3
1. Выбрать версию php
2. Выберите базу данных - mysql|mariadb|postgres|sqlite
3. Выбрать версию базы данных
4. Выберите сервер для кеширования - redis|memcached|Пропуск
5. Установка yii3

После этого в директории где выполнялась команда появится docker-compose.yml файл с настроенными сервисами.

В директории site уже будет развернут yii3 проект

## Конфигурация проекта

Конфигурация проекта - https://yiisoft.github.io/docs/guide/concept/configuration.html

Дополнительные переменные необходымые для корректной работы проекта - прокидывайте через `environment` в `docker-compose.yml`

## console

запускайте команды с помощью

```bash
docky yii {arg}
```

## Cron

По умолчанию cron включен

Если необходимо добавить задания, то сделайте публикацию файлов с заданиями

```bash
docky publish --file cron_tasks
```

Запись заданий осуществляйте в:
- `${CONF_PATH}/app/cron/docky` - для пользователя сайта
- `${CONF_PATH}/app/cron/root` - для root пользователя

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

msmtp клиент не устанавливается.
Используйте сервис mailhog для тестирования отправки почты.

- host=mailhog
- port=1025

панель доступна на - localhost:8025

## Node

Для публикации сервиса с node.js выполните команду и укажите точку входа для node.js в файле .env (переменная NODE_PATH)

```bash 
docky publish --service node
```

используйте команды ```docky npm``` чтобы запускать npm команды в контейнере

