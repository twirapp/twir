finish()
{
  docker compose -f docker-compose.dev.yml stop > /dev/null 2>&1
  echo "Services stoped, exiting."
}
trap finish EXIT SIGHUP
