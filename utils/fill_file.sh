#!/bin/bash

read -r -d '' LOG_MSG <<'EOF'
[] some log message 0\n
[] some log message 1\n
[] some log message 2\n
[] some log message 3\n
[] some log message 4\n
[] some log message 5\n
[] some log message 6\n
[] some log message 7\n
[] some log message 8\n
[] some log message 9\n
EOF

FILES_COUNT=32
MSG_BEFORE_SLEEP=512

i=0
while true; do
  if [ $i = $MSG_BEFORE_SLEEP ] ; then
    i=0
    echo "sleep"
    sleep 1
  fi
  ((i++))
  echo -e $LOG_MSG >> some_log_$(($RANDOM%32)).log
done