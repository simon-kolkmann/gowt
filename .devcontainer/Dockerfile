FROM golang:1.24.3

RUN apt update

RUN apt install -y jq nano

RUN groupadd -g 1000 dev && \
useradd -m -u 1000 -s /bin/bash -g dev dev

USER dev
