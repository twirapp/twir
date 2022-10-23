import { HttpException, Injectable } from '@nestjs/common';
import { Client, Transport } from '@nestjs/microservices';
import { config } from '@tsuwari/config';
import { AuthUser, ClientProxy } from '@tsuwari/shared';
import { In, Not } from '@tsuwari/typeorm';
import { Bot, BotType } from '@tsuwari/typeorm/entities/Bot';
import { Channel } from '@tsuwari/typeorm/entities/Channel';
import { DashboardAccess } from '@tsuwari/typeorm/entities/DashboardAccess';
import { Token } from '@tsuwari/typeorm/entities/Token';
import { User } from '@tsuwari/typeorm/entities/User';
import { AccessToken, exchangeCode, getTokenInfo } from '@twurple/auth';
import { getRawData } from '@twurple/common';
import chunk from 'lodash.chunk';

import { typeorm } from '../index.js';
import { Payload } from '../jwt/jwt.service.js';
import { JwtPayload } from '../jwt/jwt.strategy.js';
import { staticApi } from '../twitchApi.js';

@Injectable()
export class AuthService {
  @Client({ transport: Transport.NATS, options: { servers: [config.NATS_URL] } })
  nats: ClientProxy;

  async checkUser(tokens: AccessToken, userId: string, username?: string | null) {
    const defaultBot = await typeorm.getRepository(Bot).findOne({
      where: {
        type: BotType.DEFAULT,
      },
    });

    if (!defaultBot) {
      throw new Error('Bot not created, cannot create user.');
    }

    if (!tokens.refreshToken || !tokens.expiresIn) {
      throw new HttpException(
        `Something went wrong on gettings twitch tokens. Please, try again later.`,
        500,
      );
    }

    const tokenData = {
      accessToken: tokens.accessToken,
      refreshToken: tokens.refreshToken,
      obtainmentTimestamp: new Date(tokens.obtainmentTimestamp),
      expiresIn: tokens.expiresIn,
    };

    let user = await typeorm.getRepository(User).findOne({
      where: { id: userId },
      relations: {
        channel: {
          bot: true,
        },
        token: true,
      },
    });

    if (user) {
      if (!user.channel) {
        user.channel = await typeorm.getRepository(Channel).save({
          id: user.id,
          botId: defaultBot.id,
        });
      }

      if (user.tokenId) {
        await typeorm.getRepository(Token).update({ id: user.tokenId }, tokenData);
      } else {
        const token = await typeorm.getRepository(Token).save(tokenData);
        await typeorm.getRepository(User).update({ id: userId }, { token });
      }
    } else {
      const newUser = typeorm.getRepository(User).create({
        id: userId,
        token: await typeorm.getRepository(Token).save(tokenData),
      });

      user = await typeorm.manager.save(newUser);
      newUser.channel = await typeorm.getRepository(Channel).save({
        id: userId,
        botId: defaultBot.id,
      });
      await typeorm.manager.save(newUser);
    }

    if (username) {
      await this.nats
        .emit('bots.joinOrLeave', {
          action: user.channel?.isEnabled ? 'join' : 'part',
          username,
          botId: user.channel!.botId,
        })
        .toPromise();
    }

    await Promise.all([
      this.nats.send('bots.createDefaultCommands', [userId]).toPromise(),
      this.nats.emit('eventsub.subscribeToEventsByChannelId', userId).toPromise(),
    ]);

    return user;
  }

  async getProfile(userPayload: JwtPayload) {
    const [dbUser, dashboards] = await Promise.all([
      typeorm.getRepository(User).findOneBy({ id: userPayload.id }),
      typeorm.getRepository(DashboardAccess).find({
        where: { user: { id: userPayload.id } },
        relations: { channel: true },
      }),
    ]);

    if (dbUser?.isBotAdmin) {
      const channels = await typeorm.getRepository(Channel).find({
        where: {
          id: Not(In([...dashboards.map((d) => d.channelId), dbUser.id])),
        },
      });

      for (const channel of channels) {
        dashboards.push({
          id: channel.id,
          channelId: channel.id,
          userId: dbUser.id,
        });
      }
    }

    const chunks = chunk([...dashboards.map((d) => d.channelId), userPayload.id], 100);
    const twitchUsers = await Promise.all(chunks.map((c) => staticApi.users.getUsersByIds(c))).then(
      (v) => v.flat(),
    );

    const user = twitchUsers.find((u) => u.id === userPayload.id);

    if (!user || !dbUser) throw new HttpException('User not found', 404);

    const result: AuthUser = {
      ...getRawData(user),
      isTester: dbUser.isTester,
      dashboards: dashboards
        .map((d) => {
          const twitchUser = twitchUsers.find((u) => u.id === d.channelId);
          if (!twitchUser) return;
          return {
            ...d,
            twitch: getRawData(twitchUser),
          };
        })
        .filter(Boolean) as AuthUser['dashboards'],
    };

    if (dbUser.isBotAdmin) {
      result.isBotAdmin = dbUser.isBotAdmin;
    }

    return result;
  }

  async authorizeUserByTwitch(code: string, state: string): Promise<Payload> {
    const accessToken = await exchangeCode(
      config.TWITCH_CLIENTID,
      config.TWITCH_CLIENTSECRET,
      code,
      Buffer.from(state, 'base64').toString('utf-8'),
    );
    const tokenInfo = await getTokenInfo(accessToken.accessToken, config.TWITCH_CLIENTID);

    if (!tokenInfo.userId || !tokenInfo.userName) {
      throw new Error('Cannot find userId or userName in your tokenInfo');
    }

    await this.checkUser(accessToken, tokenInfo.userId, tokenInfo.userName);

    return {
      id: tokenInfo.userId,
      scopes: tokenInfo.scopes,
      username: tokenInfo.userName,
    };
  }
}
