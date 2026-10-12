#!/bin/sh
set -e

# Ensure proper permissions on the mounted SQLite data volume
mkdir -p /data
chown -R appuser:appgroup /data 2>/dev/null || true
chmod 775 /data 2>/dev/null || true

# Execute command as unprivileged appuser
exec su-exec appuser "$@"
