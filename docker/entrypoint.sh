#!/bin/sh

SCHEDULE="${CRON_SCHEDULE:-*/5 * * * *}"
echo "$SCHEDULE /usr/local/bin/ha2trmnl /config/config.yaml >> /proc/1/fd/1 2>&1" > /etc/crontabs/root

trap "echo '[entrypoint] Caught exit signal, shutting down'; exit 0" SIGINT SIGTERM

echo "[entrypoint] Running cron on schedule: $SCHEDULE"
crond -f -l 2 &
wait