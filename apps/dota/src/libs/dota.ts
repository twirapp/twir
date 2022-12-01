import { dirname, resolve } from 'path';
import { fileURLToPath } from 'url';

import { config } from '@tsuwari/config';
import { DotaGame, dotaHeroes, gameModes } from '@tsuwari/shared';
import { DotaGameMode } from '@tsuwari/typeorm/entities/DotaGameMode';
import { DotaHero } from '@tsuwari/typeorm/entities/DotaHero';
import { DotaMatch } from '@tsuwari/typeorm/entities/DotaMatch';
import protobufjs from 'protobufjs';
import SteamUser from 'steam-user';
import SteamID from 'steamid';

import { converUsers } from '../helpers/convertUsers.js';
import { redis } from './redis.js';
import { typeorm } from './typeorm.js';

export class Dota extends SteamUser {
  #watchRoot: protobufjs.Root;
  #clientHelloRoot: protobufjs.Root;
  #gcMessagesClient: protobufjs.Root;
  #gcMessagesCommon: protobufjs.Root;
  ready = false;

  sendHelloEvent() {
    if (!this.ready) return;
    const helloType = this.#clientHelloRoot.lookupType('CMsgClientHello');
    this.sendToGC(570, 4006, {}, Buffer.from(helloType.encode({}).finish()));
  }

  async init() {
    this.#watchRoot = await new protobufjs.Root().load(
      resolve(
        dirname(fileURLToPath(import.meta.url)),
        '..',
        '..',
        'protos',
        'dota2',
        'dota_gcmessages_client_watch.proto',
      ),
      {
        keepCase: true,
      },
    );

    this.#clientHelloRoot = await new protobufjs.Root().load(
      resolve(
        dirname(fileURLToPath(import.meta.url)),
        '..',
        '..',
        'protos',
        'dota2',
        'gcsdk_gcmessages.proto',
      ),
      {
        keepCase: true,
      },
    );

    this.#gcMessagesClient = await new protobufjs.Root().load(
      resolve(
        dirname(fileURLToPath(import.meta.url)),
        '..',
        '..',
        'protos',
        'dota2',
        'dota_gcmessages_client.proto',
      ),
      {
        keepCase: true,
      },
    );

    this.#gcMessagesCommon = await new protobufjs.Root().load(
      resolve(
        dirname(fileURLToPath(import.meta.url)),
        '..',
        '..',
        'protos',
        'dota2',
        'dota_gcmessages_common.proto',
      ),
      {
        keepCase: true,
      },
    );

    this.logOn({
      accountName: config.STEAM_USERNAME,
      password: config.STEAM_PASSWORD,
    });

    this.on('disconnected', (e) => {
      console.log(`Disconnected, ${e}`);
    });

    this.on('loggedOn', (e) => {
      console.log(`Logged into Steam as ${e.client_supplied_steamid} ${e.vanity_url}`);
      this.setPersona(1);
      this.gamesPlayed([570], true);
    });

    this.on('error', (e) => {
      console.log(e);
      console.error(e);
    });

    this.on('appLaunched', async () => {
      this.sendHelloEvent();
      this.ready = true;
      setInterval(() => {
        this.sendHelloEvent();
      }, 5 * 1000);
    });

    this.on('receivedFromGC', (_appId, msgId, payload) => {
      if (msgId === 8010) this.recievedFromGcCallback(payload);
    });

    return this;
  }

  async testMatchResults() {
    const type = this.#gcMessagesClient.lookupType('CMsgGCMatchDetailsRequest');
    const msg = type.encode({
      match_id: 6662322079,
    });
    this.sendToGC(570, 7095, {}, Buffer.from(msg.finish()), (_appId, msgId, payload) => {
      const type = this.#gcMessagesClient.lookupType('CMsgGCMatchDetailsResponse');
      console.log('r', type.decode(payload).toJSON());
    });
  }

  async getPresences(accs: string[]) {
    if (!this.ready) {
      return console.error('App not ready for getting presences.');
    }

    const convertedAccs = accs.map(SteamID.fromIndividualAccountID).map((id) => id.getSteamID64());
    const type = this.#watchRoot.lookupType('CMsgClientToGCFindTopSourceTVGames');

    await Promise.all([
      Promise.all(accs.map((a) => redis.del(`dotaMatches:${a}`))),
      Promise.all(accs.map((a) => redis.del(`dotaRps:${a}`))),
    ]);

    this.requestRichPresence(570, convertedAccs, 'english', async (error, data) => {
      if (error) {
        console.log(error);
        return console.error(error);
      }
      if (!data.users) return;
      const users = converUsers(data.users).filter(
        (u) => !['#DOTA_RP_INIT', '#DOTA_RP_IDLE'].includes(u.richPresence.status),
      );

      await Promise.all(
        users.map((u) =>
          redis.set(`dotaRps:${u.userId}`, JSON.stringify(u.richPresence), 'EX', 60),
        ),
      );

      const lobbyIds = new Set(
        users.filter((u) => u.richPresence.lobbyId).map((u) => u.richPresence.lobbyId),
      );
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

  getDotaProfileCard(accountId: string | number): Promise<any> {
    const type = this.#gcMessagesClient.lookupType('CMsgClientToGCGetProfileCard');
    const request = type.encode({
      account_id: Number(accountId),
    });

    const responseType = this.#gcMessagesCommon.lookupType('CMsgDOTAProfileCard');

    return new Promise((resolve) => {
      this.sendToGC(570, 7534, {}, Buffer.from(request.finish()), (_appid, msgType, payload) => {
        if (msgType === 7535) {
          const response = responseType.decode(payload).toJSON();
          console.log(response);
          if (!response.account_id) resolve(null);
          resolve(response);
        } else resolve(null);
      });
    });
  }

  async recievedFromGcCallback(payload: Buffer) {
    const type = this.#watchRoot.lookupType('CMsgGCToClientFindTopSourceTVGamesResponse');

    const data = type.decode(payload).toJSON() as {
      game_list?: Array<DotaGame>;
    };

    if (data.game_list) {
      for (const game of data.game_list) {
        if (!game.players || !game.match_id || game.players.length < 9) continue;

        const gameMode = gameModes.find((g) => g.id === game.game_mode);
        if (gameMode) {
          await typeorm
            .getRepository(DotaGameMode)
            .upsert(
              { id: game.game_mode, name: gameMode.name },
              { conflictPaths: ['id'], skipUpdateIfNoValuesChanged: true },
            );
          const gameModeEntity = await typeorm.getRepository(DotaGameMode).findOneBy({
            id: game.game_mode,
          });
          const data = {
            lobbyType: game.lobby_type,
            players: game.players.map((p) => p.account_id),
            playersHeroes: game.players.map((p) => p.hero_id),
            startedAt: new Date(Number(`${game.activate_time}000`)),
            matchId: game.match_id,
            avarageMmr: game.average_mmr,
            weekendTourneyBracketRound: game.weekend_tourney_bracket_round?.toString(),
            weekendTourneySkillLevel: game.weekend_tourney_skill_level?.toString(),
            lobbyId: game.lobby_id,
            gameMode: gameModeEntity!,
          };

          const repository = typeorm.getRepository(DotaMatch);

          if (await repository.findOneBy({ matchId: game.match_id })) {
            await repository.update({ matchId: game.match_id }, data);
          } else {
            await repository.save(data);
          }
        }

        for (const player of game.players) {
          await redis.set(`dotaMatches:${player.account_id}`, JSON.stringify(game), 'EX', 30 * 60);
          const hero = dotaHeroes.find((h) => h.id === player.hero_id);
          if (hero) {
            await typeorm
              .getRepository(DotaHero)
              .upsert(
                { id: hero.id, name: hero.localized_name },
                { conflictPaths: ['id'], skipUpdateIfNoValuesChanged: true },
              );
          }
        }
      }
    }
  }
}
