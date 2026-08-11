#!/usr/bin/env bash
set -euo pipefail

OPA_URL="http://localhost:8181/v1/data/rbac/allow"

check() {
	echo "--- $1 ---"
	curl -s -X POST "$OPA_URL" -H "Content-Type: application/json" -d "$2" | jq
}

check "admin bypasses quota" '{
  "input": {"user": {"role": "admin", "plan": "free"}, "requests_today": 999}
}'

check "free plan under limit -> allow" '{
  "input": {"user": {"role": "user", "plan": "free"}, "requests_today": 2}
}'

check "free plan at limit -> deny" '{
  "input": {"user": {"role": "user", "plan": "free"}, "requests_today": 3}
}'

check "pro plan under limit -> allow" '{
  "input": {"user": {"role": "user", "plan": "pro"}, "requests_today": 99}
}'

check "unknown plan -> deny" '{
  "input": {"user": {"role": "user", "plan": "enterprise"}, "requests_today": 0}
}'
