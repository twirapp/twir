finish()
{
  #pkill -SIGINT -F /tmp/tsuwari-redis.pid
  docker stop tsuwari-redis-stack
  docker rm tsuwari-redis-stack

  docker stop tsuwari-adminer
  docker rm tsuwari-adminer

  nats-server -sl quit=/tmp/tsuwari-nats.pid
  pg_ctl -D devbox/data/postgres stop
}
trap finish EXIT SIGHUP
