ARG HOST="localhost"
ARG VENDOR="vendor"
ARG BASE_IMAGE_VERSION="latest"

FROM ${HOST}/${VENDOR}/golang:${BASE_IMAGE_VERSION} AS builder

ARG NAME="app"

WORKDIR /tmp/build

RUN apk update \
    && apk upgrade \
    && apk add git

COPY cmd cmd
COPY internal internal
COPY pkg pkg
COPY go.* .
COPY main.go main.go

RUN mkdir -p bin \
    && go build -ldflags "-X \"main.name=${NAME}\"" -o bin/app main.go \
    && go clean -cache -modcache -testcache

FROM ${HOST}/${VENDOR}/alpine:${BASE_IMAGE_VERSION} AS final

WORKDIR /opt/app/bin

RUN apk update \
    && apk upgrade \
    && rm -rf /var/cache/apk/* /tmp/* /var/tmp/*

ENV PATH="/opt/app/bin:${PATH}"

COPY --from=builder /tmp/build/bin/app /opt/app/bin/app

COPY entrypoint.sh /

RUN ["chmod", "+x", "/entrypoint.sh"]

ENTRYPOINT ["/entrypoint.sh"]
