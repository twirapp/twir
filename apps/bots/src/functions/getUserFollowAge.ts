import { intervalToDuration, formatDuration } from 'date-fns';

import { staticApi } from '../bots.js';

export async function getuserFollowAge(fromId: string, toId: string) {
  const follow = await staticApi.users.getFollowFromUserToBroadcaster(fromId, toId);

  if (!follow) return 'not follower';

  const duration = intervalToDuration({ start: follow.followDate.getTime(), end: Date.now() });

  return formatDuration(duration);
}
