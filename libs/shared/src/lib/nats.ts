import { ClientProxy as CP } from '@nestjs/microservices';
import { Command, CommandModule, CommandPermission, Response } from '@tsuwari/prisma';
import { rawDataSymbol } from '@twurple/common';
import { EventSubChannelUpdateEvent, EventSubUserUpdateEvent } from '@twurple/eventsub';
import { Observable } from 'rxjs';

export interface ClientProxyCommands {
  'streamstatuses.process': {
    input: string[],
    result: boolean
  },
  'bots.getDefaultCommands': {
    input: any,
    result: Array<{ name: string, description?: string, visible: boolean, permission: CommandPermission, module?: CommandModule }>
  },
  'bots.getVariables': {
    input: any,
    result: Array<{
      name: string,
      example?: string,
      description?: string
    }>
  },
  'bots.deleteMessages': {
    input: {
      channelId: string,
      channelName: string,
      messageIds: string[],
    },
    result: boolean,
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
  },
  'setCommandCache': {
    input: Command & { responses?: Response[] },
    result: any,
  },
  'bots.createDefaultCommands': {
    input: string[],
    result: any,
  },
  'dota.getProfileCard': {
    input: string | number,
    result: {
      account_id: number,
      rank_tier?: number,
      leaderboard_rank?: number;
    } | null
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
  'stream.update': {
    input: EventSubChannelUpdateEvent[typeof rawDataSymbol],
    result: any,
  },
  'user.update': {
    input: EventSubUserUpdateEvent[typeof rawDataSymbol],
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
  'bots.removeTimerFromQueue': ClientProxyEvents['bots.addTimerToQueue'],
  'dota.cacheAccountsMatches': {
    input: string[],
    result: any,
  },
  'eventsub.subscribeToEventsByChannelId': {
    input: string,
    result: any,
  }
}

export type ClientProxyResult<K extends keyof ClientProxyCommands> = Observable<ClientProxyCommands[K]['result']>
export type ClientProxyCommandsKey = keyof ClientProxyCommands
export type ClientProxyEventsKey = keyof ClientProxyEvents


export abstract class ClientProxy extends CP {
  abstract send<TEvent extends keyof ClientProxyCommands>(pattern: TEvent, data: ClientProxyCommands[TEvent]['input']): Observable<ClientProxyCommands[TEvent]['result']>;
  abstract emit<TEvent extends keyof ClientProxyEvents>(pattern: TEvent, data: ClientProxyEvents[TEvent]['input']): Observable<ClientProxyEvents[TEvent]['result']>;
}
