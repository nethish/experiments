# Watchdog
Watchdog is a kernel module that does checks and acts upon. Classified in to `harddog` for hardware based watchdog and `softdog` for software based watchdog. This experiment is based on `softdog`.

Every 10s as specified in the `watchdog.conf` file, the watchdog runs tests in the `watchdog.d` directory. If there is a failure, it triggers a reboot of the system.

The `watchdog` invokes the scripts inside watchdog.d directory with `test` argument. If the script fails, then it invokes it again with `repair` argument.
