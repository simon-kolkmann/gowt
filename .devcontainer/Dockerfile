FROM golang:1.24

RUN apt update

RUN apt install -y jq nano

RUN groupadd -g 1000 dev && \
useradd -m -u 1000 -s /bin/bash -g dev dev

USER dev

RUN go install github.com/oligot/go-mod-upgrade@v0.11.0
RUN go install github.com/go-delve/delve/cmd/dlv@v1.25
