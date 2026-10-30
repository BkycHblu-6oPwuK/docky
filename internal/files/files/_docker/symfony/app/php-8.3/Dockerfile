FROM php:8.3-fpm

ARG USERGROUP
ENV SUPERVISOR_PHP_USER="docky"

RUN groupadd --force -g $USERGROUP docky
RUN useradd -ms /bin/bash --no-user-group -g $USERGROUP -u $USERGROUP docky

COPY nginx/certs/ /usr/local/share/ca-certificates/

RUN apt-get update && apt-get install -y --no-install-recommends \
    apt-utils \
    autoconf \
    pkg-config \
    build-essential \
    libssl-dev \
    libpng-dev \
    libzip-dev \
    libfreetype6-dev \
    libjpeg62-turbo-dev \
    libwebp-dev \
    libxpm-dev \
    libicu-dev \
    zip \
    unzip \
    ca-certificates

RUN apt-get update && apt-get install -y --no-install-recommends libmagickwand-dev imagemagick 

RUN apt-get update && apt-get install -y --no-install-recommends \
    wget \
    nginx \
    cron \
    supervisor \
    gifsicle \
    jpegoptim \
    optipng \
    ghostscript \
    curl \
    openssl \
    libxml2-dev \
    libonig-dev \
    libsqlite3-dev \
    libpq-dev \
    git
    
RUN set -eux; \
    apt-get update && apt-get install -y --no-install-recommends libmemcached-dev zlib1g-dev; \
    mkdir -p /tmp/pecl; cd /tmp/pecl; \
    for pkg in imagick xdebug memcached igbinary redis; do \
        curl -4 -fsSL --connect-timeout 15 --max-time 180 --retry 5 --retry-all-errors --retry-delay 3 -o "/tmp/pecl/$pkg.tgz" "https://pecl.php.net/get/$pkg"; \
    done; \
    pecl install /tmp/pecl/imagick.tgz; \
    pecl install /tmp/pecl/xdebug.tgz; \
    pecl install /tmp/pecl/memcached.tgz; \
    pecl install /tmp/pecl/igbinary.tgz; \
    pecl install -D 'enable-redis-igbinary="yes"' /tmp/pecl/redis.tgz; \
    docker-php-ext-enable imagick xdebug memcached igbinary redis; \
    rm -rf /tmp/pecl /tmp/pear

RUN docker-php-ext-configure gd --with-freetype --with-jpeg --with-webp --with-xpm && \
    docker-php-ext-install \
    pdo_mysql \
    mysqli \
    pdo_pgsql \
    pgsql \
    pdo_sqlite \
    bcmath \
    gd \
    zip \
    xml \
    intl \
    opcache


ENV COMPOSER_ALLOW_SUPERUSER=1

RUN curl -sS https://getcomposer.org/installer | php -- \
    --filename=composer \
    --install-dir=/usr/local/bin
RUN echo "alias composer='XDEBUG_MODE=off composer'" >> /etc/bash.bashrc \
 && echo "alias composer='XDEBUG_MODE=off composer'" >> /home/docky/.bashrc \
 && chown docky:docky /home/docky/.bashrc

COPY app/supervisord.conf /etc/supervisor/conf.d/supervisord.conf

COPY app/cron/ /tmp/docky-cron/
RUN set -eux; \
    mkdir -p /var/spool/cron/crontabs; \
    for u in docky root; do \
        if [ -f "/tmp/docky-cron/$u" ]; then \
            cp "/tmp/docky-cron/$u" "/var/spool/cron/crontabs/$u"; \
        else \
            touch "/var/spool/cron/crontabs/$u"; \
        fi; \
        sed -i -e 's/\r$//' -e '$a\' "/var/spool/cron/crontabs/$u"; \
        chown "$u:crontab" "/var/spool/cron/crontabs/$u"; \
        chmod 0600 "/var/spool/cron/crontabs/$u"; \
    done; \
    touch /var/log/cron.log; \
    rm -rf /tmp/docky-cron

RUN service cron restart

COPY bin/app_entrypoint.sh /usr/local/bin/entrypoint.sh
RUN sed -i 's/\r$//' /usr/local/bin/entrypoint.sh && \
    chmod +x /usr/local/bin/entrypoint.sh

COPY app/symlinks /usr/symlinks
COPY bin/create_simlink.sh /usr/local/bin/create_simlink.sh
RUN sed -i 's/\r$//' /usr/local/bin/create_simlink.sh && \
    chmod +x /usr/local/bin/create_simlink.sh

COPY app/php-8.3/php.ini /usr/local/etc/php/conf.d/php.ini
COPY app/php-8.3/xdebug.ini /usr/local/etc/php/conf.d/xdebug.ini
COPY app/php-fpm.conf /usr/local/etc/php-fpm.d/zzzzwww.conf
COPY app/nginx /etc/nginx/conf.d

RUN apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

WORKDIR /var/www

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]

