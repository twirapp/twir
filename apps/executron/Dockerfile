FROM mirror.gcr.io/debian:bookworm-slim AS prod_base
WORKDIR /app
RUN apt-get update && \
	apt-get install -y --no-install-recommends \
		apt-transport-https \
		ca-certificates \
		curl \
		gnupg \
		libstdc++6 && \
	curl -fsSL --tlsv1.2 --proto "=https" --retry 3 \
		'https://packages.doppler.com/public/cli/gpg.DE2A7741A397C129.key' | \
		gpg --dearmor -o /usr/share/keyrings/doppler-archive-keyring.gpg && \
	echo "deb [signed-by=/usr/share/keyrings/doppler-archive-keyring.gpg] https://packages.doppler.com/public/cli/deb/debian any-version main" \
		> /etc/apt/sources.list.d/doppler-cli.list && \
	apt-get update && \
	apt-get install -y --no-install-recommends doppler && \
	apt-get purge -y --auto-remove gnupg && \
	rm -rf /var/lib/apt/lists/*
COPY docker-entrypoint.sh /app/
RUN chmod +x /app/docker-entrypoint.sh
ENTRYPOINT ["/app/docker-entrypoint.sh"]

FROM prod_base
COPY apps/executron/.out/twir-executron /app/twir-executron
ENV NODE_ENV=production
CMD ["/app/twir-executron"]
