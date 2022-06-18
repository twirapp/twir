import { Logger } from '@nestjs/common';
import { config } from '@tsuwari/config';
import { ApiClient } from '@twurple/api';
import { RefreshingAuthProvider } from '@twurple/auth';
import { ChatClient, ChatSayMessageAttributes } from '@twurple/chat';
import { format } from 'date-fns';
import pc from 'picocolors';

import { increaseParsedMessages } from '../functions/increaseParsedMessages.js';
import { increaseUserMessages } from '../functions/increaseUserMessages.js';
import { ParserCache } from '../parser/cache.js';
import { ResponseParser } from '../parser/index.js';
import { CommandsParser } from './commandsParser.js';
import { GreetingsParser } from './greetingsParser.js';
import { KeywordsParser } from './keywordsParser.js';
import { ConsoleLogger } from './logger.js';
import { messageAls } from './message.als.js';
import { ModerationParser } from './moderationParser.js';
import { commandsCounter, messagesCounter } from './prometheus.js';
import { redis } from './redis.js';


export class Bot extends ChatClient {
  #api: ApiClient;
  #commandsParser: CommandsParser;
  #greetingsParser: GreetingsParser;
  #moderationParser: ModerationParser;
  #keywordsParser: KeywordsParser;

  constructor(authProvider: RefreshingAuthProvider, channels: string[]) {
    super({
      authProvider,
      channels,
      isAlwaysMod: true,
    });

    this.#commandsParser = new CommandsParser();
    this.#greetingsParser = new GreetingsParser();
    this.#moderationParser = new ModerationParser();
    this.#keywordsParser = new KeywordsParser();
    this.#api = new ApiClient({
      authProvider,
    });

    this.#registerListeners();
  }

  async say(channel: string, message: string, attributes?: ChatSayMessageAttributes) {
    const als = messageAls.getStore();
    als?.logger.log(`${pc.bgCyan(pc.black('OUT'))} ${pc.bgGreen(pc.white(channel))}: ${pc.bgYellow(pc.white(message))}`);
    if (config.isProd || config.SAY_IN_CHAT) {
      super.say(channel, message, attributes);
    }
  }

  async timeout(channel: string, user: string, duration?: number, reason?: string) {
    const isBotModRequest = await redis.get(`isBotMod:${channel.substring(1)}`);
    const isBotMod = isBotModRequest === 'true';
    if (isBotMod) {
      console.log(`${format(new Date(), 'yyyy-MM-dd\'T\'HH:mm:ss.SSSxxx')} ${pc.bgCyan(pc.black('TIMEOUT'))} ${pc.bgGreen(pc.white(channel))}: ${pc.bgYellow(pc.white(user))}`);
      super.timeout(channel, user, duration, reason);
    } else {
      console.log(`${format(new Date(), 'yyyy-MM-dd\'T\'HH:mm:ss.SSSxxx')} ${pc.bgCyan(pc.black('TIMEOUT'))} bot no mod on channel ${pc.bgGreen(pc.white(channel))}, so timeout skiped.`);
    }
  }

  async #registerListeners() {
    const me = await this.#api.users.getMe();

    this.onRegister(async () => {
      console.log(
        `${pc.bgCyan(pc.black('!'))} ${pc.magenta(me.displayName)} ${pc.green('connected to twitch servers.')}`,
      );
    });

    this.onJoin((channel) => {
      console.log(
        `${pc.bgCyan(pc.black('!'))} ${pc.magenta(me.displayName)} ${pc.green('joined a channel')} ${pc.cyan(channel.replace('#', ''))}`,
      );
    });

    this.onNamedMessage('USERSTATE', ({ tags, rawParamValues }) => {
      const channelName = rawParamValues[0]?.substring(1);
      const tag = tags.get('mod');

      if (tag === '0') {
        console.info(`${pc.bgCyan(pc.black('!'))} ${tags.get('display-name')} lost mod status in ${channelName} channel`);
        redis.del(`isBotMod:${channelName}`);
      }
      if (tag === '1') {
        console.info(`${pc.bgCyan(pc.black('!'))} ${tags.get('display-name')} got mod status in ${channelName} channel`);
        redis.set(`isBotMod:${channelName}`, 'true');
      }
    });

    this.onMessage(async (channel, user, message, state) => {
      messageAls.run({
        messageId: state.id,
        logger: new Logger(state.id),
      }, async () => {
        if (!state.channelId || !state.userInfo?.userId) return;

        const store = messageAls.getStore();
        store?.logger.log(`IN ${pc.green(channel)} | ${pc.magenta(`${user}#${state.userInfo.userId}`)}: ${pc.white(message)}`);
        const isBotModRequest = await redis.get(`isBotMod:${channel.substring(1)}`);
        const isBotMod = isBotModRequest === 'true';

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
              this.say(channel, moderateResult.message);
            }

            return;
          }
        }

        const msgObject = {
          message,
          channelId: state.channelId,
          sender: {
            id: state.userInfo.userId,
            name: state.userInfo.userName,
          },
          raw: state,
          cache: new ParserCache(state.channelId, state.userInfo.userId),
        };

        this.#commandsParser.parse(message, state).then(async (result) => {
          if (!state.channelId) return;
          if (!result?.responses?.length) return;
          commandsCounter.inc();

          for (const response of result.responses) {
            if (!response) continue;
            if (result.responses.indexOf(response) > 0 && !isBotMod) break;
            const msg = await ResponseParser.parse(response, msgObject, result.params);

            this.say(channel, msg, { replyTo: state.id });
          }
        });

        this.#greetingsParser.parse(state).then(async (response) => {
          if (!response) return;
          this.say(channel, await ResponseParser.parse(response, msgObject, message), { replyTo: state.id });
        });

        this.#keywordsParser.parse(message, state).then(async (responses) => {
          if (!responses || !responses.length) return;

          for (const response of responses) {
            if (!response) continue;
            if (responses.indexOf(response) > 0 && !isBotMod) break;
            const msg = await ResponseParser.parse(response, msgObject, message);

            this.say(channel, msg, { replyTo: state.id });
          }
        });

        increaseUserMessages(state.userInfo.userId, state.channelId);
        increaseParsedMessages(state.channelId);
        messagesCounter.inc();
      });
    });
  }
}
