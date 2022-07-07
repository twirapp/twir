import { dirname, resolve } from 'path';
import { fileURLToPath } from 'url';

import { Injectable, Logger, OnModuleInit } from '@nestjs/common';
import { config } from '@tsuwari/config';
import { PrismaClient, Prisma } from '@tsuwari/prisma';
import { DotaGame, RedisService } from '@tsuwari/shared';
import protobufjs from 'protobufjs';
import SteamUser from 'steam-user';
import SteamID from 'steamid';

import { converUsers } from './helpers/convertUsers.js';
import { dotaHeroes, gameModes } from './constants.js';

@Injectable()
export class AppService extends SteamUser implements OnModuleInit {
  #watchRoot: protobufjs.Root;
  #clientHelloRoot: protobufjs.Root;
  #logger = new Logger('Dota');

  constructor(private readonly redis: RedisService, private readonly prisma: PrismaClient) {
    super({
      autoRelogin: true,
    });
  }

  async onModuleInit() {
    this.#watchRoot = await new protobufjs.Root().load(resolve(dirname(fileURLToPath(import.meta.url)), '..', 'protos', 'dota2', 'dota_gcmessages_client_watch.proto'), {
      keepCase: true,
    });

    this.#clientHelloRoot = await new protobufjs.Root().load(resolve(dirname(fileURLToPath(import.meta.url)), '..', 'protos', 'dota2', 'gcsdk_gcmessages.proto'), {
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

    this.on('appLaunched', async (appId) => {
      const helloType = this.#clientHelloRoot.lookupType('CMsgClientHello');
      if (appId === 570) {
        this.sendToGC(570, 4006, {}, Buffer.from(helloType.encode({}).finish()));
      }
    });

    this.on('receivedFromGC', (_appId, msgId, payload) => {
      if (msgId === 8010) this.recievedFromGcCallback(payload);
    });
  }

  async getPresences(accs: string[]) {
    this.#logger.log(`Getting presences of ${accs.length} accounts.`)

    const convertedAccs = accs.map(SteamID.fromIndividualAccountID).map(id => id.getSteamID64());
    const type = this.#watchRoot.lookupType('CMsgClientToGCFindTopSourceTVGames');

    await Promise.all([
      Promise.all(accs.map(a => this.redis.del(`dotaMatches:${a}`))),
      Promise.all(accs.map(a => this.redis.del(`dotaRps:${a}`)))
    ])

    this.requestRichPresence(570, convertedAccs, 'english', async (error, data) => {
      if (error) {
        return this.#logger.error(error)
      }
      if (!data.users) return;
      const users = converUsers(data.users);

      await Promise.all(users.map(u => this.redis.set(`dotaRps:${u.userId}`, JSON.stringify(u.richPresence), 'EX', 60)))

      const lobbyIds = new Set(users.filter(u => u.richPresence.lobbyId).map(u => u.richPresence.lobbyId));
      if (!lobbyIds.size) return;

      const newMsg = type.encode({
        search_key: '',
        league_id: 0,
        hero_id: 0,
        start_game: 0,
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
        // TURBO MATCHES SHOULD BE INCLUDED, NOT SKIPPED
        if (!game.players || !game.match_id) continue;

        const gameMode = gameModes.find(g => g.id === game.game_mode)

        if (gameMode) {
          await this.prisma.dotaMatch.create({
            data: {
              lobby_type: game.lobby_type,
              players: game.players.map(p => p.account_id),
              startedAt: new Date(Number(`${game.activate_time}000`)),
              match_id: game.match_id,
              avarage_mmr: game.average_mmr,
              weekend_tourney_bracket_round: game.weekend_tourney_bracket_round,
              weekend_tourney_skill_level: game.weekend_tourney_skill_level,
              gameMode: {
                connectOrCreate: {
                  where: {
                    id: game.game_mode,
                  },
                  create: {
                    id: game.game_mode,
                    name: gameMode.name,
                  }
                }
              }
            }
          }).catch((e) => {
            if (e instanceof Prisma.PrismaClientKnownRequestError && e.code === 'P2002' && (e.meta?.target as string[]).includes('match_id')) {

            } else {
              this.#logger.log(e)
              throw new e
            }
          })
        }

        for (const player of game.players) {
          await this.redis.set(`dotaMatches:${player.account_id}`, JSON.stringify(game), 'EX', 30 * 60)
          const hero = dotaHeroes.find(h => h.id === player.hero_id)
          if (hero) {
            await this.prisma.dotaHero.create({
              data: {
                id: hero.id,
                name: hero.localized_name
              }
            }).catch((e) => {
              if (e instanceof Prisma.PrismaClientKnownRequestError && e.code === 'P2002' && (e.meta?.target as string[]).includes('id')) {

              } else throw e
            })
          }
        }
      }
    }
  }
}