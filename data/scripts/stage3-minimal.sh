# stage3-minimal.sh
#!/bin/bash

go run ./cmd \
  -vfs minimal \
  -script data/startup/stage3.txt