finish()
{
  #pkill -SIGINT -F /tmp/tsuwari-redis.pid
  docker stop tsuwari-redis-stack
  docker rm tsuwari-redis-stack
  nats-server -sl quit=/tmp/tsuwari-nats.pid
  pg_ctl -D devbox/data/postgres stop
}
trap finish EXIT SIGHUP
