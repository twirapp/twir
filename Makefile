build-dev:
	UID=${UID} GID=${GID} docker compose -f docker-compose.dev.yml build

dev:
	doppler run -- docker compose -f docker-compose.dev.yml up