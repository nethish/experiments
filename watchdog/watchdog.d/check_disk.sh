#!/bin/bash
LOG="/var/log/watchdog_test.log"
echo "$(date) - check_disk.sh running" >>"$LOG"

USAGE=$(df / | awk 'NR==2 {print $5}' | tr -d '%')
if [ "$USAGE" -gt 90 ]; then
  echo "$(date) - Disk usage high: $USAGE%" >>"$LOG"
  exit 1
fi
exit 0
