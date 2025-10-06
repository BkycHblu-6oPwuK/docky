FROM php:5.6-fpm

RUN sed -i 's|deb.debian.org|archive.debian.org|g' /etc/apt/sources.list && \
    sed -i 's|security.debian.org|archive.debian.org|g' /etc/apt/sources.list && \
    sed -i '/stretch-updates/d' /etc/apt/sources.list && \
    echo 'Acquire::Check-Valid-Until "false";' > /etc/apt/apt.conf.d/99no-check-valid-until && \
    echo 'Acquire::AllowInsecureRepositories "true";' >> /etc/apt/apt.conf.d/99no-check-valid-until

ARG USERGROUP
ENV SUPERVISOR_PHP_USER="docky"

RUN groupadd --force -g $USERGROUP docky
RUN useradd -ms /bin/bash --no-user-group -g $USERGROUP -u $USERGROUP docky

COPY nginx/certs/ /usr/local/share/ca-certificates/

RUN apt-get update && apt-get install -y --allow-unauthenticated \
    apt-utils \
    libpq-dev \
    libpng-dev \
    libzip-dev \
    libfreetype6-dev \
    libjpeg62-turbo-dev \
    libwebp-dev \
    libxpm-dev \
    zip \
    unzip

RUN apt-get install -y --allow-unauthenticated \
    jpegoptim \
    optipng \
    ghostscript \
    msmtp \
    nginx \
    ca-certificates \
    wget \
    cron \
    supervisor \
    libmagickwand-dev imagemagick \
    gifsicle \
    git

RUN update-ca-certificates && \
    pecl install imagick-3.4.4 && \
    pecl install xdebug-2.5.5 && \
    docker-php-ext-enable imagick && \
    docker-php-ext-enable xdebug && \
    docker-php-ext-configure gd --with-freetype --with-jpeg --with-webp --with-xpm && \
    docker-php-ext-install pdo_mysql mysqli bcmath gd zip

RUN pecl install memcache-3.0.8 && docker-php-ext-enable memcache

RUN pecl install igbinary-2.0.8 && \
    pecl install redis-4.3.0 && \
    docker-php-ext-enable igbinary redis

ENV COMPOSER_ALLOW_SUPERUSER=1
RUN curl -sS https://getcomposer.org/installer | php -- \
    --filename=composer \
    --install-dir=/usr/local/bin

COPY app/msmtprc /home/docky/.msmtprc
RUN chmod 600 /home/docky/.msmtprc
RUN chown docky:docky /home/docky/.msmtprc

COPY app/supervisord.conf /etc/supervisor/conf.d/supervisord.conf
COPY app/cron /var/spool/cron/crontabs
RUN echo "" >> /var/spool/cron/crontabs/docky
RUN chown docky:docky /var/spool/cron/crontabs/docky
RUN chmod 0600 /var/spool/cron/crontabs/docky /var/spool/cron/crontabs/root
RUN touch /var/log/cron.log
RUN service cron restart
COPY bin/app_entrypoint.sh /usr/local/bin/entrypoint.sh
RUN sed -i 's/\r$//' /usr/local/bin/entrypoint.sh && \
    chmod +x /usr/local/bin/entrypoint.sh
COPY app/symlinks /usr/symlinks
COPY bin/create_simlink.sh /usr/local/bin/create_simlink.sh
RUN sed -i 's/\r$//' /usr/local/bin/create_simlink.sh && \
    chmod +x /usr/local/bin/create_simlink.sh

COPY app/php-7.4/php.ini /usr/local/etc/php/conf.d/php.ini
COPY app/php-7.4/xdebug.ini /usr/local/etc/php/conf.d/xdebug.ini
COPY app/php-fpm.conf /usr/local/etc/php-fpm.d/zzzzwww.conf
COPY app/nginx /etc/nginx/conf.d

RUN apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

WORKDIR /var/www

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
