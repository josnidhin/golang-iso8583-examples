FROM golang:1.19.1-bullseye as builder

ENV BUILDER_APP_DIR=/root/golang-iso8583-examples/example-3

WORKDIR ${BUILDER_APP_DIR}
COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN make release-build

## Build the production image
FROM alpine:3.15

ENV BUILDER_APP_DIR=/root/golang-iso8583-examples/example-3
ENV APP_DIR=/app/srv/golang-iso8583-examples/example-3
ENV APP_NAME=example-3

RUN apk update \
  && adduser -D -h /home/golang golang

WORKDIR ${APP_DIR}

COPY --from=builder ${BUILDER_APP_DIR}/${APP_NAME} .

RUN chown -R golang:golang $APP_DIR

USER golang

#CMD ${APP_DIR}/${APP_NAME}
