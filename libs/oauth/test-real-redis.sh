#!/usr/bin/env bash
set -euo pipefail

container="twir-oauth-test-${RANDOM}-$$"
cleanup() {
	docker rm -f "${container}" >/dev/null 2>&1 || true
}
trap cleanup EXIT INT TERM

docker run --detach --name "${container}" --publish 127.0.0.1::6379 redis:7.4-alpine >/dev/null
mapping="$(docker port "${container}" 6379/tcp)"
port="${mapping##*:}"
for _ in {1..100}; do
	if docker exec "${container}" redis-cli ping >/dev/null 2>&1; then
		break
	fi
	sleep 0.05
done
docker exec "${container}" redis-cli ping >/dev/null

TWIR_OAUTH_TEST_REDIS_ADDR="127.0.0.1:${port}" \
TWIR_OAUTH_TEST_REDIS_CONTAINER="${container}" \
"$@"
