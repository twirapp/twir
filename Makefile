build-dev:
	UID=${UID} GID=${GID} docker compose -f docker-compose.dev.yml build

up-dev:
	doppler run -- docker compose -f docker-compose.dev.yml up -d

down-dev:
	docker compose -f docker-compose.dev.yml down