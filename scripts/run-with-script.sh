#!/bin/bash

go run ./cmd \
  --vfs simple \
  --prompt "user@localhost:~$ " \
  --script startup.txt