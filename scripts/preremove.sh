#!/bin/sh
set -e
# Stop the service if running as a systemd unit
if command -v systemctl >/dev/null 2>&1; then
    if systemctl is-active --quiet servdocs 2>/dev/null; then
        echo "Stopping servdocs ..."
        systemctl stop servdocs || true
        systemctl disable servdocs || true
    fi
fi
