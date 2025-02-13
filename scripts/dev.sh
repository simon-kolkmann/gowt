#!/bin/bash

# This script runs the app with live reloading for development purposes.
#
# See also:
#   - watchexec: https://github.com/watchexec/watchexec
#   - discussion: https://github.com/charmbracelet/bubbletea/issues/150#issuecomment-2492038498
watchexec -r -c --wrap-process session -e go -- "DEBUG=true go run ."
