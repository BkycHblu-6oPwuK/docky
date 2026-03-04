# Конфигурация под nuxt в битрикс

php 8.2, 8.3, 8.4, 8.5

точка входа для npm - директория nuxt в корне сайта (site/nuxt)

для режима разработки выполните команду `docky npm run dev`

настройка сервера

```
if ($request_uri ~ "^(.+\.[a-zA-Z0-9]{2,5})/$") {
    return 301 $1;
}

location ~* (/\.ht|/\.git|/\.gitignore|/\.bash|/\.sql|\.settings\.php|/composer|/bitrix/cache|/bitrix/backup|/bitrix/updates|/bitrix/php_interface|/bitrix/stack_cache|/bitrix/managed_cache|/bitrix/html_pages/\.|/upload/1c_exchange|local/modules|local/php_interface|/logs/) { 
    deny all; 
}

# php админка
location /bitrix {
    try_files $uri $uri/ @bitrix;
    add_header Last-Modified $date_gmt;
    add_header Cache-Control 'no-store, no-cache';
    if_modified_since off;
    expires off;
    etag off;
}
# php api
location /api {
    try_files $uri $uri/ @bitrix;
}
# php api
location /local {
    try_files $uri $uri/ @bitrix;
}
# php файлы
location ~ \.php$ {
   try_files $uri @bitrix;
   fastcgi_split_path_info ^(.+\.php)(/.+)$;
   fastcgi_pass app:9000;
   fastcgi_index index.php;
   include fastcgi_params;
   fastcgi_param SCRIPT_FILENAME $document_root$fastcgi_script_name;
   fastcgi_param PATH_INFO $fastcgi_path_info;
   fastcgi_param SERVER_NAME $host;
   fastcgi_param SERVER_PORT $server_port;
   fastcgi_param SERVER_PROTOCOL $server_protocol;
   fastcgi_read_timeout 300s;
   fastcgi_connect_timeout 300s;
   fastcgi_send_timeout 300s;
}
# 404 для битрикс
location @bitrix {
    fastcgi_pass app:9000;
    include fastcgi_params;
    fastcgi_param SERVER_NAME $host;
    fastcgi_param SERVER_PORT $server_port;
    fastcgi_param SCRIPT_FILENAME /var/www/bitrix/urlrewrite.php;
    fastcgi_read_timeout 300s;
    fastcgi_connect_timeout 300s;
    fastcgi_send_timeout 300s;
}
# статика
location ^~ /upload {
    access_log off;
    expires 7d;
    add_header Cache-Control public;
    add_header X-Content-Type-Options nosniff;
    add_header Access-Control-Allow-Origin "*";
    try_files $uri =404;
}

location ~* \.(?:css(\.map)?|js(\.map)?|jpe?g|png|gif|ico|cur|heic|webp|tiff?|mp3|m4a|aac|ogg|midi?|wav|mp4|mov|webm|mpe?g|avi|ogv|flv|wmv)$ {
    expires    7d;
    access_log off;
    add_header Cache-Control public;
    add_header X-Content-Type-Options nosniff;
}

location ~* \.(?:svgz?|ttf|ttc|otf|eot|woff2?)$ {
    add_header Access-Control-Allow-Origin "*";
    expires    7d;
    access_log off;
    add_header Cache-Control public;
    add_header X-Content-Type-Options nosniff;
}
# сервер nuxt
location ^~ /_nuxt/ {
    proxy_pass http://node:5174/_nuxt/;
    proxy_http_version 1.1;

    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;

    proxy_set_header Upgrade $http_upgrade;
    proxy_set_header Connection $connection_upgrade;

    access_log off;
    expires 7d;
    add_header Cache-Control public;
}

location ^~ /__vite_ping {
    proxy_pass http://node:5174/__vite_ping/;
}

location ^~ /__nuxt_devtools__ {
    proxy_pass http://node:5174;
}

location / {
    proxy_pass http://node:5174;
    proxy_http_version 1.1;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header Upgrade $http_upgrade;
    proxy_set_header Connection $connection_upgrade;
}
```

# pm2

команда pm2 в docker-compose.yml для сервиса node

```yml
command: sh -c "pm2 start /var/www/nuxt/ecosystem.config.cjs --name node-server & tail -f /dev/null"
```