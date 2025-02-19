#!/bin/bash

curl -L -o /tmp/watchexec.deb https://github.com/watchexec/watchexec/releases/download/v2.3.0/watchexec-2.3.0-x86_64-unknown-linux-gnu.deb
apt install /tmp/watchexec.deb
