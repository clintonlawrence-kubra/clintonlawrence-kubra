#!/usr/bin/env bash
set -euo pipefail

# Apps to close before shutdown. Edit this list as needed.
APPS=(
  "Google Chrome"
  "Slack"
  "Visual Studio Code"
  "Finder"
  "Gmail"
  "Microsoft Teams"
  "TextMate"
)

for app in "${APPS[@]}"; do
  echo "Quitting $app..."
  # "quit saving yes" auto-saves unsaved docs instead of prompting; apps that
  # don't support the verb (e.g. Finder) just ignore it and quit normally.
  osascript -e "tell application \"$app\" to quit saving yes" 2>/dev/null || true
done

# give apps a chance to actually exit
sleep 5

# force-quit anything still open so a stray dialog can't block shutdown
for app in "${APPS[@]}"; do
  if pgrep -f "$app" >/dev/null 2>&1; then
    echo "Force-quitting $app..."
    killall "$app" 2>/dev/null || true
  fi
done

echo "Shutting down..."
sudo shutdown -h now
