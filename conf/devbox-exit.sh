finish()
{
  kill $(jobs -p)
}
trap finish EXIT SIGHUP
