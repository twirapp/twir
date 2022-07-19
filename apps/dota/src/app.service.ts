import { dirname, resolve } from 'path';
import { fileURLToPath } from 'url';

import { Injectable, Logger, OnModuleInit } from '@nestjs/common';
import { config } from '@tsuwari/config';
import { PrismaClient, Prisma } from '@tsuwari/prisma';
import { DotaGame, dotaHeroes, gameModes, RedisService } from '@tsuwari/shared';
import protobufjs from 'protobufjs';
import SteamUser from 'steam-user';
import SteamID from 'steamid';

import { converUsers } from './helpers/convertUsers.js';

@Injectable()
export class AppService extends SteamUser implements OnModuleInit {
  #watchRoot: protobufjs.Root;
  #clientHelloRoot: protobufjs.Root;
  #matchRoot: protobufjs.Root;
  #logger = new Logger('Dota');

  constructor(private readonly redis: RedisService, private readonly prisma: PrismaClient) {
    super({
      autoRelogin: true,
    });
  }


  sendHelloEvent() {
    const helloType = this.#clientHelloRoot.lookupType('CMsgClientHello');
    this.sendToGC(570, 4006, {}, Buffer.from(helloType.encode({}).finish()));
    this.#logger.log('Sent hello event.');
  }

  async onModuleInit() {
    this.#watchRoot = await new protobufjs.Root().load(resolve(dirname(fileURLToPath(import.meta.url)), '..', 'protos', 'dota2', 'dota_gcmessages_client_watch.proto'), {
      keepCase: true,
    });

    this.#clientHelloRoot = await new protobufjs.Root().load(resolve(dirname(fileURLToPath(import.meta.url)), '..', 'protos', 'dota2', 'gcsdk_gcmessages.proto'), {
      keepCase: true,
    });

    this.#matchRoot = await new protobufjs.Root().load(resolve(dirname(fileURLToPath(import.meta.url)), '..', 'protos', 'dota2', 'dota_gcmessages_client.proto'), {
      keepCase: true,
    });

    this.logOn({
      accountName: config.STEAM_USERNAME,
      password: config.STEAM_PASSWORD,
    });

    this.on('loggedOn', (e) => {
      this.#logger.log(`Logged into Steam as ${e.client_supplied_steamid} ${e.vanity_url}`);
      this.setPersona(1);
      this.gamesPlayed([570], true);
    });

    this.on('error', (e) => {
      console.log(e);
      this.#logger.error(e);
    });

    this.on('appLaunched', async (appId) => {
      this.sendHelloEvent();

      setInterval(() => {
        this.sendHelloEvent();
      }, 5 * 1000);
    });

    this.on('receivedFromGC', (_appId, msgId, payload) => {
      if (msgId === 8010) this.recievedFromGcCallback(payload);
    });
  }

  async testMatchResults() {
    const type = this.#matchRoot.lookupType('CMsgGCMatchDetailsRequest');
    const msg = type.encode({
      match_id: 6662322079,
    });
    this.sendToGC(570, 7095, {}, Buffer.from(msg.finish()), (_appId, msgId, payload) => {
      const type = this.#matchRoot.lookupType('CMsgGCMatchDetailsResponse');
      console.log('r', type.decode(payload).toJSON());
    });
  }

  async getPresences(accs: string[]) {
    this.#logger.log(`Getting presences of ${accs.length} accounts.`);

    const convertedAccs = accs.map(SteamID.fromIndividualAccountID).map(id => id.getSteamID64());
    const type = this.#watchRoot.lookupType('CMsgClientToGCFindTopSourceTVGames');

    await Promise.all([
      Promise.all(accs.map(a => this.redis.del(`dotaMatches:${a}`))),
      Promise.all(accs.map(a => this.redis.del(`dotaRps:${a}`))),
    ]);

    this.requestRichPresence(570, convertedAccs, 'english', async (error, data) => {
      if (error) {
        console.log(error);
        return this.#logger.error(error);
      }
      if (!data.users) return;
      const users = converUsers(data.users).filter(u => !['#DOTA_RP_INIT', '#DOTA_RP_IDLE'].includes(u.richPresence.status));

      await Promise.all(users.map(u => this.redis.set(`dotaRps:${u.userId}`, JSON.stringify(u.richPresence), 'EX', 60)));

      const lobbyIds = new Set(users.filter(u => u.richPresence.lobbyId).map(u => u.richPresence.lobbyId));
      if (!lobbyIds.size) return;

      const newMsg = type.encode({
        search_key: '',
        league_id: 0,
        hero_id: 0,
        start_game: 90,
        game_list_index: 0,
        lobby_ids: [...lobbyIds.values()],
      });
      this.sendToGC(570, 8009, {}, Buffer.from(newMsg.finish()));
    });
  }


  async recievedFromGcCallback(payload: Buffer) {
    const type = this.#watchRoot.lookupType('CMsgGCToClientFindTopSourceTVGamesResponse');

    const data = type.decode(payload).toJSON() as {
      game_list?: Array<DotaGame>
    };

    if (data.game_list) {
      for (const game of data.game_list) {
        if (!game.players || !game.match_id || game.players.length < 9) continue;

        const gameMode = gameModes.find(g => g.id === game.game_mode);

        if (gameMode) {
          const data = {
            lobby_type: game.lobby_type,
            players: game.players.map(p => p.account_id),
            players_heroes: game.players.map(p => p.hero_id),
            startedAt: new Date(Number(`${game.activate_time}000`)),
            match_id: game.match_id,
            avarage_mmr: game.average_mmr,
            weekend_tourney_bracket_round: game.weekend_tourney_bracket_round?.toString(),
            weekend_tourney_skill_level: game.weekend_tourney_skill_level?.toString(),
            lobbyId: game.lobby_id,
            gameMode: {
              connectOrCreate: {
                where: {
                  id: game.game_mode,
                },
                create: {
                  id: game.game_mode,
                  name: gameMode.name,
                },
              },
            },
          };

          if (await this.prisma.dotaMatch.findFirst({ where: { match_id: game.match_id } })) {
            await this.prisma.dotaMatch.update({
              where: { match_id: game.match_id },
              data,
            });
          } else {
            await this.prisma.dotaMatch.create({ data });
          }
        }

        for (const player of game.players) {
          await this.redis.set(`dotaMatches:${player.account_id}`, JSON.stringify(game), 'EX', 30 * 60);
          const hero = dotaHeroes.find(h => h.id === player.hero_id);
          if (hero) {
            await this.prisma.dotaHero.create({
              data: {
                id: hero.id,
                name: hero.localized_name,
              },
            }).catch((e) => {
              if (e instanceof Prisma.PrismaClientKnownRequestError && e.code === 'P2002' && (e.meta?.target as string[]).includes('id')) {
                null;
              } else throw e;
            });
          }
        }
      }
    }
  }
}