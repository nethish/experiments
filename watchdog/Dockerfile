FROM debian:stable-slim

RUN apt-get update && apt-get install -y watchdog && \
    rm -rf /var/lib/apt/lists/*

COPY watchdog.conf /etc/watchdog.conf
COPY watchdog.d/ /etc/watchdog.d/

CMD ["/usr/sbin/watchdog", "-F"]   
# -F keeps it in foreground for Docker logs
