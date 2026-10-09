#!/bin/bash

go run ./cmd \
  --vfs deep \
  --prompt "user@localhost" \
  --script data/startup/commands-test.txt