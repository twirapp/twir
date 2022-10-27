import { Logger } from '@nestjs/common';
import { config } from '@tsuwari/config';
import * as Parser from '@tsuwari/nats/parser';
import { Channel } from '@tsuwari/typeorm/entities/Channel';
import { ApiClient } from '@twurple/api';
import { RefreshingAuthProvider } from '@twurple/auth';
import { ChatClient, ChatSayMessageAttributes } from '@twurple/chat';
import { format } from 'date-fns';
import pc from 'picocolors';

import { increaseParsedMessages } from '../functions/increaseParsedMessages.js';
import { increaseUserMessages } from '../functions/increaseUserMessages.js';
import { storeUserMessage } from '../functions/storeUserMessage.js';
import { type ResolvableValue } from '../typings/index.js';
import { GreetingsParser } from './greetingsParser.js';
import { KeywordsParser } from './keywordsParser.js';
import { ModerationParser } from './moderationParser.js';
import { nats } from './nats.js';
import {
  commandsCounter,
  commandsResponseTime,
  greetingsCounter,
  greetingsParseTime,
  keywordsCounter,
  keywordsParseTime,
  messageParseTime,
  messagesCounter,
  moderationParseTime,
} from './prometheus.js';
import { typeorm } from './typeorm.js';

const channelRepository = typeorm.getRepository(Channel);
const strRegexp = /.{1,500}(\s|$)/;
export class Bot extends ChatClient {
  #api: ApiClient;
  #greetingsParser: GreetingsParser;
  #moderationParser: ModerationParser;
  #keywordsParser: KeywordsParser;
  #logger: Logger;

  constructor(
    authProvider: RefreshingAuthProvider,
    channels: ResolvableValue<string[]>,
    botId: string,
  ) {
    super({
      authProvider,
      channels,
      isAlwaysMod: true,
    });

    this.#greetingsParser = new GreetingsParser();
    this.#moderationParser = new ModerationParser();
    this.#keywordsParser = new KeywordsParser();
    this.#api = new ApiClient({
      authProvider,
    });

    this.#registerListeners();
    this.#logger = new Logger(botId);
  }

  async say(channel: string, message: string, attributes?: ChatSayMessageAttributes) {
    this.#logger.log(
      `${pc.bgCyan(pc.black('OUT'))} ${pc.bgGreen(pc.white(channel))}: ${pc.bgYellow(
        pc.white(message),
      )}`,
    );
    if (config.isProd || config.SAY_IN_CHAT) {
      for (const str of message.match(strRegexp)!.map((s) => s.trim())) {
        super.say(channel, str, attributes);
      }
    }
  }

  async myTimeout(
    channelName: string,
    channelId: string,
    user: string,
    duration?: number,
    reason?: string,
  ) {
    const isBotModRequest = await channelRepository.findOneBy({
      id: channelId,
    });
    const isBotMod = isBotModRequest?.isBotMod;
    if (isBotMod) {
      console.log(
        `${format(new Date(), `yyyy-MM-dd'T'HH:mm:ss.SSSxxx`)} ${pc.bgCyan(
          pc.black('TIMEOUT'),
        )} ${pc.bgGreen(pc.white(channelName))}: ${pc.bgYellow(pc.white(user))}`,
      );
      this.timeout(channelName, user, duration, reason);
    } else {
      console.log(
        `${format(new Date(), `yyyy-MM-dd'T'HH:mm:ss.SSSxxx`)} ${pc.bgCyan(
          pc.black('TIMEOUT'),
        )} bot no mod on channel ${pc.bgGreen(pc.white(channelName))}, so timeout skiped.`,
      );
    }
  }

  async #registerListeners() {
    const me = await this.#api.users.getMe();

    this.onRegister(async () => {
      console.log(
        `${pc.bgCyan(pc.black('!'))} ${pc.magenta(me.displayName)} ${pc.green(
          'connected to twitch servers.',
        )}`,
      );
    });

    this.onJoin((channel) => {
      console.log(
        `${pc.bgCyan(pc.black('!'))} ${pc.magenta(me.displayName)} ${pc.green(
          'joined a channel',
        )} ${pc.cyan(channel.replace('#', ''))}`,
      );
    });

    this.onMessage(async (channel, user, message, state) => {
      if (!state.channelId || !state.userInfo?.userId) return;
      const perfStart = performance.now();
      messagesCounter.inc();

      const replyTo = state.tags.get('reply-parent-display-name');
      if (replyTo) {
        message = message.replace(`@${replyTo}`, '').trim();
      }

      const isBotModRequest = await channelRepository.findOneBy({ id: state.channelId });
      const isBotMod = isBotModRequest?.isBotMod;

      const isModerate = !state.userInfo.isBroadcaster && !state.userInfo.isMod && isBotMod;
      if (isModerate) {
        const moderateResult = await this.#moderationParser.parse(message, state);

        if (moderateResult) {
          if (moderateResult.delete) {
            this.deleteMessage(channel, state.id);
          } else {
            this.timeout(channel, user, moderateResult.time, moderateResult.message ?? undefined);
          }

          if (moderateResult.message) {
            await this.say(channel, moderateResult.message);
          }

          moderationParseTime.observe(performance.now() - perfStart);
          return;
        }
      }

      const usersBadges: string[] = [];
      if (state.userInfo.isBroadcaster) usersBadges.push('BROADCASTER');
      if (state.userInfo.isMod) usersBadges.push('MODERATOR');
      if (state.userInfo.isSubscriber || state.userInfo.isFounder) usersBadges.push('SUBSCRIBER');
      if (state.userInfo.isVip) usersBadges.push('VIP');
      usersBadges.push('VIEWER');

      if (message.startsWith('!')) {
        const data = Parser.Request.toBinary({
          channel: {
            id: state.channelId,
            name: state.target.value.replace('#', ''),
          },
          message: {
            id: state.id,
            text: message,
          },
          sender: {
            badges: usersBadges,
            displayName: state.userInfo.displayName,
            id: state.userInfo.userId,
            name: state.userInfo.userName,
          },
        });
        nats
          .request('parser.handleProcessCommand', data, {
            timeout: 5 * 5000,
          })
          .then(async (r) => {
            const { responses: result, isReply, keepOrder } = Parser.Response.fromBinary(r.data);
            commandsCounter.inc();

            for (const response of result) {
              if (!response) continue;
              if (result.indexOf(response) > 0 && !isBotMod) break;

              const args = { replyTo: isReply ? state.id : undefined };
              if (keepOrder) {
                await this.say(channel, response, args);
              } else {
                this.say(channel, response, args);
              }
            }

            commandsResponseTime
              .labels(channel, message.split(' ')[0] ?? '')
              .observe(performance.now() - perfStart);
          });
      }

      this.#greetingsParser.parse(state).then(async (response) => {
        if (!response) return;
        const data = Parser.ParseResponseRequest.toBinary({
          channel: {
            id: state.channelId!,
            name: state.target.value.replace('#', ''),
          },
          message: {
            id: state.id,
            text: response.text,
          },
          sender: {
            badges: usersBadges,
            displayName: state.userInfo.displayName,
            id: state.userInfo.userId,
            name: state.userInfo.userName,
          },
        });

        const request = await nats.request('parser.parseTextResponse', data);
        const responseData = Parser.ParseResponseResponse.fromBinary(request.data);

        if (responseData) {
          for (const r of responseData.responses) {
            await this.say(channel, r, { replyTo: response.isReply ? state.id : undefined });
          }
          greetingsCounter.inc();
          greetingsParseTime.observe(performance.now() - perfStart);
        }
      });

      this.#keywordsParser.parse(message, state).then(async (responses) => {
        if (!responses || !responses.length) return;

        for (const response of responses) {
          if (!response) continue;
          if (responses.indexOf(response) > 0 && !isBotMod) break;
          const data = Parser.ParseResponseRequest.toBinary({
            channel: {
              id: state.channelId!,
              name: state.target.value.replace('#', ''),
            },
            message: {
              id: state.id,
              text: response.response,
            },
            sender: {
              badges: usersBadges,
              displayName: state.userInfo.displayName,
              id: state.userInfo.userId,
              name: state.userInfo.userName,
            },
          });

          const request = await nats.request('parser.parseTextResponse', data);
          const responseData = Parser.ParseResponseResponse.fromBinary(request.data);

          if (responseData) {
            for (const r of responseData.responses) {
              this.say(channel, r, { replyTo: response.isReply ? state.id : undefined });
            }
            keywordsCounter.inc();
            keywordsParseTime.observe(performance.now() - perfStart);
          }
        }
      });

      increaseUserMessages(state.userInfo.userId, state.channelId);
      increaseParsedMessages(state.channelId);
      storeUserMessage(state, message);
      messageParseTime.observe(performance.now() - perfStart);
    });
  }
}
