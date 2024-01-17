FROM golang:alpine AS build

ENV CGO_ENABLED=0

RUN apk add --no-cache \
      git \
      curl \
      make

WORKDIR /src
COPY . .
RUN make install

FROM ubuntu:bionic
LABEL maintainer="Marc Wickenden <marc@4armed.com>"

RUN \
  apt-get clean &&\
  apt-get update -y &&\
  apt-get install --no-install-recommends -y \
  # chromium dependencies
  libnss3 \
  libxss1 \
  libasound2 \
  libxtst6 \
  libgtk-3-0 \
  libgbm1 \
  ca-certificates \
  # fonts
  fonts-liberation fonts-noto-color-emoji fonts-noto-cjk \
  # timezone
  tzdata \
  # headful mode support, for example: $ xvfb-run chromium-browser --remote-debugging-port=9222
  xvfb \
  # cleanup
  && apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

RUN useradd -u 1001 -m -s /bin/bash kubeshot

COPY --from=build /go/bin/kubeshot /usr/local/bin/kubeshot

WORKDIR /data

RUN chown kubeshot:kubeshot /data

USER kubeshot
