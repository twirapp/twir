import { dirname, resolve } from 'path';
import { fileURLToPath } from 'url';

import { Injectable, Logger, OnModuleInit } from '@nestjs/common';
import { config } from '@tsuwari/config';
import { PrismaClient } from '@tsuwari/prisma';
import { RedisService } from '@tsuwari/shared';
import protobufjs from 'protobufjs';
import SteamUser from 'steam-user';
import SteamID from 'steamid';

import { converUsers } from './helpers/convertUsers.js';
import { Game } from './types.js';

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

  getPresences(accs: string[]) {
    this.#logger.log(`Getting presences of ${accs.length} accounts.`)

    const convertedAccs = accs.map(SteamID.fromIndividualAccountID).map(id => id.getSteamID64());
    const type = this.#watchRoot.lookupType('CMsgClientToGCFindTopSourceTVGames');

    this.requestRichPresence(570, convertedAccs, 'english', (error, data) => {
      if (error) {
        return this.#logger.error(error)
      }
      if (!data.users) return;
      const users = converUsers(data.users);

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
      game_list?: Array<Game>
    };

    if (data.game_list) {
      for (const game of data.game_list) {
        if (!game.average_mmr || !game.players || !game.match_id) continue;

        for (const player of game.players) {
          await this.redis.set(`dotaMatches:${player.account_id}`, JSON.stringify(game), 'EX', 30 * 60)
        }
      }
    }
  }
}