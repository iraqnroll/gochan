#!/bin/sh
set -e

if [ "$(stat -c %u /app/static/content)" != "$(id -u gochan)" ]; then
  chown -R gochan:gochan /app/static/content
fi

exec "$@"