import { CommandPermission } from '@tsuwari/prisma';
import type { TwitchPrivateMessage } from '@twurple/chat/lib/commands/TwitchPrivateMessage.js';

export type DefaultCommand = {
  name: string,
  description?: string,
  visible?: boolean,
  example?: string,
  permission: CommandPermission,
  handler: (state: TwitchPrivateMessage, params?: string) => undefined | string | string[] | Promise<string[] | string | undefined> | Promise<undefined | string>
}