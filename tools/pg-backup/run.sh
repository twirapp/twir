#! /bin/sh

set -e


if [ "${BACKUP_SCHEDULE}" = "" ]; then
  sh backup.sh
else
  echo "$BACKUP_SCHEDULE sh /app/backup.sh" > /etc/crontabs/root
  sh backup.sh
  exec crond -l 2 -f
fi