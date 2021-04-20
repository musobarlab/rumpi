# Builder
# ---------------------------------------------------
FROM golang:1.14-alpine as builder

WORKDIR /usr/src/app

COPY . .

RUN go clean -modcache

RUN go mod download \
    && CGO_ENABLED=0 GOOS=linux \
        go build -ldflags '-w -s' -a -o rumpi github.com/musobarlab/rumpi/cmd/rumpi
# ---------------------------------------------------


# Final Image
# ---------------------------------------------------
FROM alpine

RUN apk update && \
    apk add --no-cache tzdata

ARG DIST_ENVIRONMENT=development

ENV TZ=Asia/Jakarta

WORKDIR /usr/src/app
ENV APP_PATH=/usr/src/app

# copy all configuration file to APP_PATH
COPY config/key/app.key ${APP_PATH}/config/key/
COPY config/key/app.key.pub ${APP_PATH}/config/key/

# COPY .env.${DIST_ENVIRONMENT} ${APP_PATH}/.env
COPY .env ${APP_PATH}/.env

COPY --from=builder /usr/src/app/rumpi ${APP_PATH}/

EXPOSE 9001

ENTRYPOINT ["/usr/src/app/rumpi"]
# ---------------------------------------------------
