FROM golang:1.24.3

RUN apt update

RUN apt install -y jq nano

RUN curl -L -o /tmp/watchexec.deb https://github.com/watchexec/watchexec/releases/download/v2.3.0/watchexec-2.3.0-x86_64-unknown-linux-gnu.deb
RUN apt install /tmp/watchexec.deb

RUN groupadd -g 1000 dev && \
useradd -m -u 1000 -s /bin/bash -g dev dev

USER dev
