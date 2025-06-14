# Обертка под docker compose

Настроенная обертка над docker compose для локальных проектов Bitrix, Laravel.

- Под Bitrix смотрите в [bitrix.md](docs/bitrix.md)
- Под Laravel смотрите в [laravel.md](docs/laravel.md)
- Под Symfony смотрите в [symfony.md](docs/symfony.md)
- Под ванильный php смотрите в [vanilla.md](docs/vanilla.md)

А здесь общее описания работы с скриптом ```docky```

## Установка

Для установки можно запустить установочный скрипт с правами суперпользователя:

```bash
curl -sSL https://raw.githubusercontent.com/BkycHblu-6oPwuK/docky/main/scripts/install.sh | sudo sh
```
Либо же можете вручую скачать бинарник (находится в bin/docky) и поместить его в ```/usr/local/bin``` и не забудьте дать ему необходимые права (команда ```chmod +x```)

Обновление скрипта происходит точно также, при выполнении команды curl - файл заменяется.

После установки проверьте работу скрипта, можете выполнить команду ```docky```

## Публикация docker-compose.yml

Команды выполняются в директории с docker-compose.yml или в любой другой дочерней, но при этом файл docker-compose.yml должен существовать (для всех команд, кроме - init и clean-cache).

Для публикации docker-compose.yml выполните команду:

```bash
docky init
```

Если docker-compose.yml уже существует, то будет предложено создать новый - в таком случае файл будет переименован и после этого вы перейдете к дальшейшей публикации нового файла, а иначе будет произведен выход.
Сам проект размещается в директории ```site```, которая создается в той же директории где находиться docker-compose.yml

## Конфигурация yml через файл .env

Файл автоматически создается при инициализации проекта. В него вы так же можете добавлять свои переменные и менять стандартные

```
PHP_VERSION=7.4|8.2|8.3|8.4 # по сути не фактическая версия php - используется для построения пути к dockerfile в ```_docker/app/php-{PHP_VERSION}/dockerfile```, в самом докерфайле задается фактическая версия
MYSQL_VERSION={любая версия доступная на docker hub}
MARIADB_VERSION={любая версия доступная на docker hub}
POSTGRES_VERSION={любая версия доступная на docker hub}
NODE_VERSION={любая версия доступная по ссылке - https://deb.nodesource.com/setup_${NODE_VERSION}.x}
NODE_PATH=/var/www/local/js/vite # здесь путь до package json в контейнере, поэтому указывайте вместе абсолютный путь. /var/www - это DOCUMENT_ROOT сайта в контейнере
USERGROUP={id группы пользователя (обычно 1000), по умолчанию скрипт автоматически прокидывает id группы текущего пользователя консоли, но если вы запустите скрипт из под root то автоматически пробросит 1000. Используйте эту переменную если нужно изменить группу пользователя в контейнерах app,nginx,node}
```

Переменные окружения прокидываемые скриптом:

```
SITE_PATH - путь к директории site в директории с docker-compose.yml
DOCKER_PATH - путь к директории _docker в ~.cache/docky/_docker, либо _docker, если такая существует, в одной директории с docker-compose.yml
CONF_PATH - путь к директории _conf в одной директории с docker-compose.yml
USERGROUP - id группы текущего пользователя консоли
```

## Публикация докерфайлов и файлов конфигурации

Если вам необходимо внести изменения в докерфайлы или файлы конфигурации, или добавить что то свое, то используйте команду:

```bash
docky publish
```

При этом если директория ./_docker уже существует, то она будет переименована.

Можно опубликовать отдельные файлы командой:

```bash
docky publish --file php.ini|xdebug.ini|cron_tasks|nginx_conf|symlinks|mysql_conf|postgres_conf|supervisord_conf
docky publish --dockerfile app|nginx|node
```

Публикация происходит в директорию ```${CONF_PATH}```
- php.ini|xdebug.ini|app-dockerfile - ${CONF_PATH}/app/php-${PHP_VERSION}/*
- cron_tasks - ${CONF_PATH}/app/cron
- symlinks - ${CONF_PATH}/app/symlinks
- mysql_conf - ${CONF_PATH}/mysql/my.cnf для сервисов mysql и mariadb
- postgres_conf - ${CONF_PATH}/postgres/postgres.conf
- supervisord_conf - ${CONF_PATH}/app/supervisord.conf
- nginx_conf - ${CONF_PATH}/nginx/conf.d
- nginx_dockerfile - ${CONF_PATH}/nginx/Dockerfile
- node_dockerfile - ${CONF_PATH}/node/Dockerfile

Так же добавляются необходимые volumes в сервисы docker-compose.yml

опубликовать отдельный сервис в docker-compose.yaml:

```bash
docky publish --service node|mysql|mariadb|postgres|sphinx|redis|memcached|mailhog|phpmyadmin
```

## nginx

Сертификаты и ключи копируются в контейнер из /_docker/nginx/certs/ и запись о них уже добавлена в nginx.conf.

Размещайте свои сертификаты в ${CONF_PATH}/nginx/certs/ и добавляйте каждый сертификат через volumes в docker-compose.yml

```
- ${CONF_PATH}/nginx/certs/site:/usr/local/share/ca-certificates/site
- ${CONF_PATH}/nginx/certs/cert.crt:/usr/local/share/ca-certificates/cert.crt
```

Сейчас там используются временные самописные сертификаты которые будут действительны до ~2051 года.

Сервер одинаково настроен на работу как по http, так и по https.

Более подробно - [certificates.md](docs/certificates.md)

Выполните шаги из пункта "Импорт в windows" чтобы не было ошибок в браузере

Так же вы можете добавлять файлы конфигурации для nginx через volumes в docker-compose.yml

```
- ${CONF_PATH}/nginx/site.conf:/etc/nginx/conf.d/site.conf
```

для работы сокетов на локальном сайте в php контейнере так же устанавливается nginx который проксирует запросы на основной контейнер с nginx.

## php

Конфигурации для каждых из версий находятся по пути - ``` ${DOCKER_PATH}/app/php-{PHP_VERSION}/ ```.

Для изменения версии - измените ее в файле .env, переменая ```PHP_VERSION```

свои конфигурации можно размещать в - ``` ${CONF_PATH}/app/php-{PHP_VERSION}/ ```.

php.ini публикуется командой:

```bash
docky publish --file php.ini
```

файл будет помещен в - ``` ${CONF_PATH}/app/php-{PHP_VERSION}/php.ini ```

## Xdebug

По умолчанию установлен.

xdebug.ini публикуется командой:

```bash
docky publish --file xdebug.ini
```

файл будет помещен в - ``` ${CONF_PATH}/app/php-{PHP_VERSION}/xdebug.ini ```

## Node и npm, npx

По умолчанию установлена 23 версия.

Так же по умолчанию контейнер работает на двух портах - ``` 5173 ``` и ``` 5174 ```

При установке можно указать корневую директории для вашего фронтенда, то есть та директория, где расположен ``` package.json ```. Контейнер работает внутри данной директории и соответственно там выполняются команды.

Вы можете изменить версию ``` node ```  и корневую директорию, поменяв значение переменных ``` NODE_VERSION ``` и ``` NODE_PATH ``` в файле .env

Скачивание node идет с адреса - https://deb.nodesource.com/setup_${NODE_VERSION}.x, и часто подверсий там нет, указывайте целые числа или можете перейти по ссылке и проверить доступность версии.

Команды npm/npx выполняются слудующим образом:
```bash
docky npm {arg}
docky npx {arg}
```

для работы сервера разработки vite обавьте это в vite.config.js
```js
    server: {
        host: '0.0.0.0',
        port: 5173,
        open: false,
        cors: {
            origin: '*'
        },
        hmr: {
            host: 'localhost',
        },
    }
```

## Пакет pm2 в node контейнере

В dockerfile node устанавливается пакет pm2, для работы с ним используйте команду

```bash 
docky pm2 {arg}
```

Запустить сервер node js можно:

1. командой pm2 - ```docky pm2 start server.js```
2. настроить на запуск при запуске контейнеров - для этого добавьте команду (в сервисе node) - ``` command: sh -c "pm2 start /var/www/local/js/vite/server.js --name node-server & tail -f /dev/null" ```, здесь указывается точка входа к скрипту который будет запускаться

## Туннелирование локального сайта

Для туннелирования используется Expose (https://github.com/beyondcode/expose) и для того чтобы поделиться вашим локальным сайтом выполните команду:

```bash
docky share
```

Сайт будет доступен 1 час, после этого команду можно выполнить заново.

привер флагов expose для проброса через команду docky share:

- `--auth`
- `--server`
- `--subdomain`
- `--domain`
- `--server-host`
- `--server-port`
- ...

Документация expose - https://expose.dev/docs/introduction

## Символические ссылки

При запуске контейнера app запускается скрипт - _docker/bin/create_simlink.sh, он создает ссылки внутри контейнера и соответственно ссылки внутри сайта распространяются и на хост и другие контейнеры.

Ссылки берутся из файла - ${DOCKER_PATH}/app/symlinks. Структура файла должна быть такой:

```
/var/www/<path> /var/www/<path>
```

Если же вам нужно дополнительные ссылки добавить, выполните команду:

```bash
docky publish --file symlinks
```

файл появится в - ${CONF_PATH}/app/symlinks

формируйте все пути относительно структуры контейнера.


## Создание нового домена для основного сайта

0. Выполните команду ```docky down``` или убедитесь что контейнеры остановлены
1. Для создания выполните команду ``` docky create domain ``` или ``` docky create domain --name {domain_name} ```
2. Введите доменное имя сайта
3. Выполните команду ``` docky build ```

## Добавление записей в hosts 

При создании домена создается файл hosts в ``` ${CONF_PATH}/hosts ```
А так же необходимые конфиги для nginx и сертификаты для этого домена в 
- ${CONF_PATH}/app/nginx/${domain}.conf
- ${CONF_PATH}/nginx/conf.d/${domain}.conf
- ${CONF_PATH}/nginx/conf.d/snippets/${domain}.conf
- ${CONF_PATH}/nginx/certs/${domain}/*

и необходимые volumes для этих файлов в docker-compose.yml

В целом вы можете добавлять в него записи вида:

```
127.0.0.1 new_site
```

После этого можно выполнить команду:

```bash
docky hosts push
```

Но лучше выполнять команду ``` docky create domain ```, если вы хотите чтобы создались сертификаты и конфиги для nginx

И все записи из вашего локального hosts будут добавлены в глобальный (${SYSTEM_DISK}\Windows\System32\drivers\etc\hosts - если wsl, или /etc/hosts - если ubuntu)

(${SYSTEM_DISK} - скрипт через команду powershell попытается найти системный диск)

Если записи уже существуют, то дублирования не будет.

## Описание всех доступных команд

- `version|v` - версия программы
```bash
docky version
docky v
```
- `update` - обновление программы
```bash
docky update
```
- `init` - Создание docker-compose.yml в текущей директории и создает директорию ```site```
```bash
docky init
```
- `publish` - Публикация файлов конфигурации в директории с docker-compose.yml, доступен флаг ``` service ``` для публикации отдельного сервиса в docker-compose.yml
```bash
docky publish
docky publish --service node|mysql|mariadb|postgres|sphinx|redis|memcached|mailhog|phpmyadmin
docky publish --file php.ini|xdebug.ini|cron_tasks|nginx_conf|symlinks|mysql_conf|postgres_conf|supervisord_conf
docky publish --dockerfile app|nginx|node
```
- `clean-cache` - очищает кэш директории скрипта, в ней храняться файлы конфигурации, докерфайлы
```bash
docky clean-cache
```
- `create site` - Создание нового сайта в директории сайта (./site/new-site.ru) для фреймворка bitrix
```bash
docky create site
docky create site --name {domain_name}
```
- `create domain` - Создание нового домена для вашего основного сайта
```bash
docky create domain
docky create domain --name {domain_name}
```
- `hosts push` - Переносит записи из вашего локального hosts файла в глобальный
```bash
docky hosts push
```
- `share` - Позволяет сделать сайт доступным из интернета
```bash
docky share
```
- `php` - Выполнение команды php в контейнере с php
```bash
docky php -v
```
- `artisan` - Выполнение команды artisan в контейнере с php для фреймворка laravel
```bash
docky artisan {arg}
```
- `symfony` - Выполнение команды bin/console в контейнере с php для фреймворка symfony
```bash
docky symfony {arg}
```
- `composer` - Выполнение команды composer в контейнере с php
```bash
docky composer install
```
- `npm` - Выполнение команды npm в контейнере с node js если он был установлен
```bash
docky npm install
```
- `npx` - Выполнение команды npx в контейнере с node js если он был установлен
```bash
docky npx create-vite@latest my-app
```
- `pm2` - Выполнение команды pm2 в контейнере с node js если он был установлен. Команда принимает такие же аргументы как и оригинальная pm2
```bash
docky pm2 {arg}
```
- `reset` - сбрасывает docker-compose.yml под актуальную версию скрипта, при этом старый файл будет переименован и вы всегда можете откатить изменение
```bash
docky reset
```
- `И все дефолтные команды docker-compose` - Выполнение любой команды docker-compose
```bash
docky up -d
docky down
docky build
```

## Пользователи в контейнерах

- `docky` - в контейнере с php (service app), nginx, node

## Настройка Redis

Опубликуйте сервис командой

```bash
docky publish --service redis
```

## Настройка memcached

Опубликуйте сервис командой

```bash
docky publish --service memcached
```