import { CommandConditional } from '../libs/commandsParser.js';
import { redis } from '../libs/redis.js';

export async function getChannelCommandsNamesFromRedis(channelId: string) {
  const channelCommandsKeys = await redis.keys(`commands:${channelId}:*`);

  if (!channelCommandsKeys.length) return;

  const channelCommandsNames = channelCommandsKeys.map((c) => c.split(':')[2]) as string[];
  if (!channelCommandsNames || !channelCommandsNames.length) return;

  return channelCommandsNames;
}

export async function getChannelCommandsByNamesFromRedis(channelId: string, names: string[]) {
  const result: CommandConditional[] = [];

  for (const name of names) {
    const command: CommandConditional = (await redis.hgetall(
      `commands:${channelId}:${name}`,
    )) as unknown as CommandConditional;

    if (!Object.keys(command).length) continue;
    if ((JSON.parse(command.aliases as string) as Array<string>).includes(name)) continue;
    result.push(command);
  }

  return result;
}
