import type { TwitchPrivateMessage } from '@twurple/chat/lib/commands/TwitchPrivateMessage.js';

export type DefaultCommand = {
  name: string,
  description?: string,
  visible?: string,
  example?: string,
  handler: (state: TwitchPrivateMessage, params?: string) => undefined | string | string[] | Promise<string[] | undefined> | Promise<undefined | string>
}