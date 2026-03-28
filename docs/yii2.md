# Конфигурация под yii2

nginx + php (8.2, 8.3, 8.4, 8.5) + mysql|mariadb|postgres|sqlite
а так же redis|memcached

## Шаги публикации docker-compose.yml

Шаги:

0. Выбор фреймворка - yii2
1. Выберите шаблон - advanced|basic
2. Выбрать версию php
3. Выберите базу данных - mysql|mariadb|postgres|sqlite
4. Выбрать версию базы данных
5. Выберите сервер для кеширования - redis|memcached|Пропуск
6. Установка yii2

После этого в директории где выполнялась команда появится docker-compose.yml файл с настроенными сервисами.

В директории site уже будет развернут yii2 проект

Если выбирали шаблон advanced, то не забудте выполнить запуск скрипта init

```bash
docky php /var/www/init
```

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

