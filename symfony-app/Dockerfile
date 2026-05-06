FROM php:8.4-fpm-alpine

# 1. Install system dependencies
RUN apk add --no-cache \
    libpng-dev \
    libjpeg-turbo-dev \
    freetype-dev \
    rabbitmq-c-dev \
    autoconf \
    build-base \
    libzip-dev \
    zip \
    unzip

# 2. Install PHP extensions
RUN docker-php-ext-install gd zip
RUN pecl install amqp && docker-php-ext-enable amqp

# 3. Set the working directory
WORKDIR /var/www/html

# 4. Install Composer
COPY --from=composer:latest /usr/bin/composer /usr/bin/composer
COPY php.ini /usr/local/etc/php/conf.d/docker-php-config.ini

# 5. Copy your files into the container
COPY . .

# 6. Set permissions
RUN chown -R www-data:www-data /var/www/html