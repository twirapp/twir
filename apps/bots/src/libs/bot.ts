import { Logger } from '@nestjs/common';
import { config } from '@tsuwari/config';
import { ApiClient } from '@twurple/api';
import { RefreshingAuthProvider } from '@twurple/auth';
import { ChatClient, ChatSayMessageAttributes } from '@twurple/chat';
import { format } from 'date-fns';
import pc from 'picocolors';

import { increaseParsedMessages } from '../functions/increaseParsedMessages.js';
import { increaseUserMessages } from '../functions/increaseUserMessages.js';
import { nestApp } from '../nest/index.js';
import { ParserService } from '../nest/parser/parser.service.js';
import { GreetingsParser } from './greetingsParser.js';
import { KeywordsParser } from './keywordsParser.js';
import { ModerationParser } from './moderationParser.js';
import { commandsCounter, commandsResponseTime, greetingsCounter, greetingsParseTime, keywordsCounter, keywordsParseTime, messageParseTime, messagesCounter, moderationParseTime } from './prometheus.js';
import { redis } from './redis.js';

export class Bot extends ChatClient {
  #api: ApiClient;
  #greetingsParser: GreetingsParser;
  #moderationParser: ModerationParser;
  #keywordsParser: KeywordsParser;
  #parserService = nestApp.get(ParserService);
  #logger: Logger;

  constructor(authProvider: RefreshingAuthProvider, channels: string[], botId: string) {
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
    this.#logger.log(`${pc.bgCyan(pc.black('OUT'))} ${pc.bgGreen(pc.white(channel))}: ${pc.bgYellow(pc.white(message))}`);
    if (config.isProd || config.SAY_IN_CHAT) {
      for (const str of message.match(/.{1,500}(\s|$)/g)!.map(s => s.trim())) {
        super.say(channel, str, attributes);
      }
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
      if (!state.channelId || !state.userInfo?.userId) return;
      const perfStart = performance.now();
      messagesCounter.inc();

      this.#logger.log(`IN ${pc.green(channel)} | ${pc.magenta(`${user}#${state.userInfo.userId}`)}: ${pc.white(message)}`);
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
            await this.say(channel, moderateResult.message);
          }

          moderationParseTime.observe(performance.now() - perfStart);
          return;
        }
      }

      if (message.startsWith('!')) {
        this.#parserService.parseChatMessage(state.rawLine!).then(async (result) => {
          if (!state.channelId) return;
          if (!result?.length) return;
          commandsCounter.inc();

          for (const response of result) {
            if (!response) continue;
            if (result.indexOf(response) > 0 && !isBotMod) break;

            await this.say(channel, response, { replyTo: state.id });
          }

          commandsResponseTime.labels(channel, message.split(' ')[0] ?? '').observe(performance.now() - perfStart);
        });
      }

      this.#greetingsParser.parse(state).then(async (response) => {
        if (!response) return;
        const result = await this.#parserService.parseResponse({
          channelId: state.channelId!,
          userId: state.userInfo.userId,
          userName: state.userInfo.userName,
          text: response,
        });

        if (result) {
          for (const r of result) {
            await this.say(channel, r, { replyTo: state.id });
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
          const result = await this.#parserService.parseResponse({
            channelId: state.channelId!,
            userId: state.userInfo.userId,
            userName: state.userInfo.userName,
            text: response,
          });

          if (result) {
            for (const r of result) {
              await this.say(channel, r, { replyTo: state.id });
            }
            keywordsCounter.inc();
            keywordsParseTime.observe(performance.now() - perfStart);
          }
        }
      });

      increaseUserMessages(state.userInfo.userId, state.channelId);
      increaseParsedMessages(state.channelId);
      messageParseTime.observe(performance.now() - perfStart);
    });
  }
}
