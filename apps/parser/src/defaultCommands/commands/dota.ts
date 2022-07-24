import { ID } from '@node-steam/id';
import { config } from '@tsuwari/config';
import { Prisma, PrismaService } from '@tsuwari/prisma';
import { ClientProxy, dotaHeroes, dotaMedals, gameModes, RedisService } from '@tsuwari/shared';
import { HelixStreamData } from '@twurple/api/lib/index.js';
import axios from 'axios';
import rateLimit from 'axios-rate-limit';

import { app } from '../../index.js';
import { DefaultCommand } from '../types.js';

const prisma = app.get(PrismaService);
const redis = app.get(RedisService);
const nats = app.get('NATS').providers[0].useValue as ClientProxy;

const messages = Object.freeze({
  GAME_NOT_FOUND: 'Game not found.' as string,
  NO_ACCOUNTS: 'You have not added account.' as string,
});

const dotaApiInstance = axios.create({
  baseURL: `http://api.steampowered.com/`,
  timeout: 1000,
});
dotaApiInstance.interceptors.request.use((req) => {
  req.params = req.params || {};
  req.params['key'] = config.STEAM_API_KEY;
  return req;
});
const dotaApi = rateLimit(dotaApiInstance, { maxRequests: 2, perMilliseconds: 1000, maxRPS: 2 });


// DO NOT CHANGE ORDER!
const colors = ['Blue', 'Teal', 'Purple', 'Yellow', 'Orange', 'Pink', 'Gray', 'Light Blue', 'Green', 'Brown'];

const getPlayerHero = (heroId: number, index?: number) => {
  if (heroId === 0 && typeof index !== 'undefined') {
    const color = colors[index];
    return color ?? 'Unknown';
  } else if (heroId === 0 && typeof index === 'undefined') return 'Unknown';
  else {
    const hero = dotaHeroes.find(h => h.id === heroId);
    if (!hero) return 'Unknown';
    return hero.localized_name;
  }
};

const getGames = async (accounts: string[], take = 1) => {
  const rps = await Promise.all(accounts.map(a => redis.get(`dotaRps:${a}`)));
  if (!rps.filter(r => r !== null).length) {
    return messages.GAME_NOT_FOUND;
  }

  const cachedGames = await Promise.all(accounts.map(a => redis.get(`dotaMatches:${a}`)));
  if (!cachedGames.filter(r => r !== null).length) {
    return messages.GAME_NOT_FOUND;
  }

  const dbGames = await prisma.dotaMatch.findMany({
    where: {
      players: {
        hasSome: accounts.map(a => Number(a)),
      },
    },
    orderBy: {
      startedAt: 'desc',
    },
    include: {
      gameMode: true,
      playersCards: true,
    },
    take,
  });

  if (!dbGames.length) return messages.GAME_NOT_FOUND;

  return dbGames.map(g => {
    const players = g.players.map((p, index) => ({ account_id: p, hero_id: g.players_heroes[index]! }));

    return {
      ...g,
      players,
    };
  });
};

const getAccounts = async (channelId: string) => {
  const accounts = await prisma.dotaAccount.findMany({
    where: {
      channelId,
    },
  });

  return accounts.length ? accounts : messages.NO_ACCOUNTS;
};

export const dota: DefaultCommand[] = [
  {
    name: 'dota addacc',
    permission: 'BROADCASTER',
    visible: false,
    async handler(state, params?) {
      if (!params || !state.channelId) return;

      const id = new ID(params).getAccountID();
      if (Number.isNaN(id)) return 'Wrong account id.';

      try {
        await prisma.dotaAccount.create({
          data: {
            channelId: state.channelId,
            id: params,
          },
        });

      } catch (e) {
        if (e instanceof Prisma.PrismaClientKnownRequestError && e.code === 'P2002' && (e.meta?.target as string[]).includes('id')) {
          return `Account ${id} already added.`;
        } else throw e;
      }

      return 'Account added.';
    },
  },
  {
    name: 'dota delacc',
    permission: 'BROADCASTER',
    visible: false,
    async handler(state, params?) {
      if (!params || !state.channelId) return;

      const id = new ID(params).getAccountID();
      if (Number.isNaN(id)) return 'Wrong account id.';

      const account = await prisma.dotaAccount.findUnique({
        where: {
          id_channelId: {
            channelId: state.channelId,
            id: id.toString(),
          },
        },
      });

      if (!account) return `Account ${id} not linked.`;

      await prisma.dotaAccount.delete({
        where: {
          id_channelId: {
            channelId: state.channelId,
            id: params,
          },
        },
      });

      return 'Account deleted.';
    },
  },
  {
    name: 'np',
    permission: 'VIEWER',
    visible: true,
    handler: async (state) => {
      if (!state.channelId) return;

      const accounts = await getAccounts(state.channelId);
      if (typeof accounts === 'string') return accounts;

      const games = await getGames(accounts.map(a => a.id), 1);
      if (typeof games === 'string') return games;

      return games
        .map(g => {
          const avgMmr = g.gameMode.id === 22 && g.lobby_type === 7 ? ` (${g.avarage_mmr}mmr)` : '';
          return `${g.gameMode.name}${avgMmr}`;
        })
        .join(' | ');
    },
  },
  {
    name: 'wl',
    permission: 'VIEWER',
    visible: true,
    handler: async (state) => {
      if (!state.channelId) return;

      const stream = await redis.get(`streams:${state.channelId}`);

      if (!stream) {
        return 'Stream is offline';
      }

      const parsedStream = JSON.parse(stream) as HelixStreamData;
      const accounts = await getAccounts(state.channelId);
      if (typeof accounts === 'string') return accounts;

      const games = await prisma.dotaMatch.findMany({
        where: {
          startedAt: {
            gte: new Date(new Date(parsedStream.started_at).getTime() - 10 * 60 * 1000),
          },
          players: {
            hasSome: accounts.map(a => Number(a.id)),
          },
          lobby_type: {
            in: [0, 7],
          },
        },
        orderBy: {
          startedAt: 'desc',
        },
        select: {
          match_id: true,
          gameMode: true,
          players_heroes: true,
          players: true,
          finished: true,
          result: true,
        },
      }).then(gms => gms.filter(g => !g.players_heroes.some(h => h === 0)));

      const gamesForRequest = games.filter(g => !g.result);
      const matchesData: any[] = [...games.filter(g => g.result).map(g => g.result)];
      const createResults: any[] = [];

      for (const game of gamesForRequest) {
        try {
          const request = await dotaApi.get('IDOTA2Match_570/GetMatchDetails/v1', { params: { match_id: game.match_id } });
          if (request.status !== 200) continue;
          const data = request.data;
          if (data.result.error) continue;
          if (data.result) {
            matchesData.push(data.result);
            createResults.push(data.result);
          }
        } catch (e) {
          console.error(e);
          continue;
        }
      }

      if (createResults.length) {
        const data = createResults.map(d => ({
          match_id: d.match_id.toString(),
          players: d.players,
          radiant_win: d.radiant_win,
          game_mode: d.game_mode,
        }));

        await prisma.dotaMatchResult.createMany({
          data,
        });
      }

      const matchesByGameMode: {
        [x: number]: {
          matches: any[],
          stringResult: string
        }
      } = {};
      gameModes.forEach(m => {
        matchesByGameMode[m.id] = {
          matches: [],
          stringResult: '',
        };
      });

      for (const account of accounts) {
        for (const match of matchesData) {
          if (typeof match.radiant_win === 'undefined' || !match.players) continue;

          await prisma.dotaMatch.update({
            where: {
              match_id: match.match_id.toString(),
            },
            data: {
              finished: true,
            },
          });

          let player = match.players.find((p: any) => p.account_id === Number(account.id));
          if (!player) {
            const dbMatch = games.find(g => g.match_id === match.match_id.toString());
            if (!dbMatch) continue;
            const playerIndex = dbMatch.players.indexOf(Number(account.id));
            player = match.players[playerIndex];
          }

          if (!player) continue;

          const hero = dotaHeroes.find(h => h.id === player.hero_id);
          const isPlayerRadiant = player.team_number === 0;
          let isWinner: boolean;

          if ((isPlayerRadiant && match.radiant_win) || (!isPlayerRadiant && !match.radiant_win)) {
            isWinner = true;
          } else {
            isWinner = false;
          }

          matchesByGameMode[match.game_mode]?.matches.push({
            isWinner,
            hero,
            kills: player.kills,
            deaths: player.deaths,
            assists: player.assists,
          });
        }
      }

      const result: string[] = [];

      for (const [modeId, data] of Object.entries(matchesByGameMode).filter(e => e[1].matches.length)) {
        const wins = data.matches.filter(r => r.isWinner);
        const mode = gameModes.find(m => m.id === Number(modeId));
        const heroesResult = data.matches
          .filter(m => typeof m.hero !== 'undefined')
          .map(m => `${m.hero.localized_name}(${m.isWinner ? 'W' : 'L'}) [${m.kills}/${m.deaths}/${m.assists}]`)
          .reverse();

        let msg = `${mode?.name ?? 'Unknown'} W ${wins.length} — L ${data.matches.length - wins.length}`;

        if (heroesResult.join(', ').length > 500) {
          const heroesResultShort = data.matches
            .filter(m => typeof m.hero !== 'undefined')
            .map(m => `${m.hero.shortName ?? m.hero.localized_name}(${m.isWinner ? 'W' : 'L'}) [${m.kills}/${m.deaths}/${m.assists}]`)
            .reverse()
            .join(', ');

          if (heroesResultShort.length <= 500) {
            msg += `: ${heroesResultShort}`;
          } else {
            msg += `: ${data.matches
              .filter(m => typeof m.hero !== 'undefined')
              .map(m => `${m.hero.shortName ?? m.hero.localized_name}(${m.isWinner ? 'W' : 'L'})`)
              .reverse().join(', ')}`;
          }
        } else msg += `: ${heroesResult.join(', ')}`;
        matchesByGameMode[Number(modeId)]!.stringResult = msg;
        result.push(msg);
      }
      //send from matchesByGameMode
      return result.length ? result : 'W 0 — L 0';
    },
  },
  {
    name: 'dota listacc',
    permission: 'BROADCASTER',
    visible: false,
    async handler(state) {
      if (!state.channelId) return;

      const accounts = await getAccounts(state.channelId);
      return typeof accounts === 'string' ? accounts : accounts.map(a => a.id).join(', ');
    },
  },
  {
    name: 'lg',
    permission: 'VIEWER',
    visible: false,
    async handler(state) {
      if (!state.channelId) return;

      const accounts = await getAccounts(state.channelId);
      if (typeof accounts === 'string') return accounts;

      const accountsIds = accounts.map(a => a.id);
      const games = await getGames(accountsIds);
      if (typeof games === 'string') return games;

      const prevGame = games[1];
      if (!prevGame) return 'Previous game not found.';

      const currentGame = games[0]!;

      const neededPlayers: Array<{
        prev: Player,
        curr: Player,
      }> = [];

      for (const player of currentGame.players.filter(p => !accountsIds.includes(p.account_id.toString()))) {
        const findedPlayer = prevGame.players.find(p => p.account_id === player.account_id);

        if (!findedPlayer) continue;
        neededPlayers.push({
          prev: findedPlayer,
          curr: player,
        });
      }

      if (!neededPlayers.length) return 'Not playing with anyone from last game.';
      return neededPlayers.map((p) => `${getPlayerHero(p.curr.hero_id)} played as ${getPlayerHero(p.prev.hero_id)}`).join(', ');
    },
  },
  {
    name: 'gm',
    permission: 'VIEWER',
    visible: false,
    async handler(state) {
      if (!state.channelId) return;

      const accounts = await getAccounts(state.channelId);
      if (typeof accounts === 'string') return accounts;

      const accountsIds = accounts.map(a => a.id);
      const games = await getGames(accountsIds);
      if (typeof games === 'string') return games;

      const game = games[0]!;

      const usersForGet = game.players.filter(p => !game.playersCards.some(c => c.account_id === p.account_id.toString()));

      const cardsResponses = await Promise.all(usersForGet.map((u) =>
        nats.send('dota.getProfileCard', u.account_id).toPromise(),
      ));

      cardsResponses.forEach(async (c) => {
        if (!c) return;
        await prisma.dotaMatchProfileCard.create({
          data: {
            account_id: c.account_id.toString(),
            match_id: game.id,
            rank_tier: c.rank_tier,
            leaderboard_rank: c.leaderboard_rank,
          },
        });
      });

      const users = [
        ...game.playersCards.map(p => ({
          account_id: p.account_id,
          rank_tier: p.rank_tier,
          leaderboard_rank: p.leaderboard_rank,
        })),
        ...cardsResponses.map(p => ({
          account_id: p!.account_id.toString(),
          rank_tier: p!.rank_tier,
          leaderboard_rank: p?.leaderboard_rank,
        })),
      ];

      const result: string[] = [];

      for (const p of users) {
        const playerIndex = game.players.map(p => p.account_id).indexOf(Number(p.account_id));
        const player = game.players[playerIndex]!;
        const medal = dotaMedals.find(m => m.rank_tier === p.rank_tier) || { rank_tier: 0, name: 'Unknown' };
        const rank = medal.rank_tier === 80 && p.leaderboard_rank ? `#${p.leaderboard_rank}` : '';

        result[playerIndex] = `${getPlayerHero(player.hero_id, playerIndex)}: ${medal.name}${rank}`;
      }

      return result.join(', ');
    },
  },
];

type Player = {
  account_id: number;
  hero_id: number;
}