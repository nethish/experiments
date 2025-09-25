#!/bin/bash
LOG="/var/log/watchdog_test.log"
echo "$(date) - check_failure.sh running" >>"$LOG"

# Force a failure for testing
if [ ! -f /tmp/fail ]; then
  touch /tmp/fail
fi
exit 1
