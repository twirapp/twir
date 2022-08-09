import { ApiClient, HelixUserData } from '@twurple/api';
import { getRawData } from '@twurple/common';
import Redis from 'ioredis';
import _ from 'lodash';

export async function getTwitchUsers(
  input: string,
  redis: Redis,
  api: ApiClient,
): Promise<HelixUserData>;
export async function getTwitchUsers(
  input: string[],
  redis: Redis,
  api: ApiClient,
): Promise<HelixUserData[]>;
export async function getTwitchUsers(
  input: string | string[],
  redis: Redis,
  api: ApiClient,
): Promise<HelixUserData | HelixUserData[]> {
  if (Array.isArray(input)) {
    const users: Array<HelixUserData> = [];
    const chunks = _.chunk(input, 100);

    for (const chunk of chunks) {
      const cachedUsers = await Promise.all(
        input.map((id) => redis.get(`twitchCachedUsers:${id}`)),
      ).then((users) => {
        return users
          .map((u) => (u ? (JSON.parse(u) as HelixUserData) : null))
          .filter((u) => u !== null);
      });
      const cachedUsersIds = new Set(cachedUsers.map((u) => u.id));
      const usersForGet = chunk.filter((id) => !cachedUsersIds.has(id));

      const twitchUsers = await api.users
        .getUsersByIds(usersForGet)
        .then((users) => users.map(getRawData));
      twitchUsers.forEach((u) =>
        redis.set(`twitchCachedUsers:${u.id}`, JSON.stringify(u), 'EX', 5 * 60 * 60),
      );

      users.push(...cachedUsers, ...twitchUsers);
    }

    return users;
  } else {
    const user = await redis.get(`twitchCachedUsers:${input}`);
    if (!user) {
      const twitchUser = await api.users.getUserById(input).then((u) => getRawData(u));
      redis.set(
        `twitchCachedUsers:${twitchUser.id}`,
        JSON.stringify(twitchUser),
        'EX',
        5 * 60 * 60,
      );

      return twitchUser;
    } else {
      return JSON.parse(user) as HelixUserData;
    }
  }
}
