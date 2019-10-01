FROM ruby:2.6.4-alpine3.9

WORKDIR /app

ENV RUN_DEPS="mariadb-client mariadb-dev" \
    DEV_DEPS="make gcc musl-dev"

COPY Gemfile Gemfile.lock /app/

RUN apk add --no-cache --update --virtual .rundeps ${RUN_DEPS} &&\
    apk add --no-cache --update --virtual .devdeps ${DEV_DEPS} &&\
    gem install bundler -v 2.0.2 &&\
    bundle install &&\
    apk del .devdeps

COPY . .

ENTRYPOINT [ "rake" ]
