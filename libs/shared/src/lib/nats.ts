import { ClientProxy as CP } from '@nestjs/microservices';
import { CommandPermission } from '@tsuwari/prisma';
import { Observable } from 'rxjs';

export interface ClientProxyCommands {
  'streamstatuses.process': {
    input: string[],
    result: boolean
  },
  'bots.getDefaultCommands': {
    input: any,
    result: Array<{ name: string, description?: string, permission: CommandPermission }>
  },
  'bots.getVariables': {
    input: any,
    result: Array<{
      name: string,
      example?: string,
      description?: string
    }>
  },
  'parseChatMessage': {
    input: string,
    result: string[]
  },
  'parseResponse': {
    input: {
      userId?: string,
      channelId: string,
      userName?: string,
      userDisplayName?: string,
      text: string
    };
    result: string[];
  }
}

export interface ClientProxyEvents {
  'streams.online': {
    input: { streamId: string, channelId: string },
    result: any
  },
  'streams.offline': {
    input: { channelId: string },
    result: any
  },
  'bots.joinOrLeave': {
    input: {
      action: 'join' | 'part',
      username: string,
      botId: string,
    },
    result: any
  },
  'bots.addTimerToQueue': {
    input: string,
    result: any
  },
  'bots.removeTimerFromQueue': ClientProxyEvents['bots.addTimerToQueue']
}

export type ClientProxyResult<K extends keyof ClientProxyCommands> = Observable<ClientProxyCommands[K]['result']>

export abstract class ClientProxy extends CP {
  abstract send<TEvent extends keyof ClientProxyCommands>(pattern: TEvent, data: ClientProxyCommands[TEvent]['input']): Observable<ClientProxyCommands[TEvent]['result']>;
  abstract emit<TEvent extends keyof ClientProxyEvents>(pattern: TEvent, data: ClientProxyEvents[TEvent]['input']): Observable<ClientProxyEvents[TEvent]['result']>;
}
