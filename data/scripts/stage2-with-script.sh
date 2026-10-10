# stage2-with-script.sh
#!/bin/bash

go run ./cmd \
  -vfs deep \
  -prompt "vedsatt@plasma" \
  -script data/startup/stage2.txt