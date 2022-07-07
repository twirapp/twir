import { config } from '@tsuwari/config';
import { PrismaService } from '@tsuwari/prisma';
import { DotaGame, dotaHeroes, gameModes, RedisService, TwitchApiService } from '@tsuwari/shared';
import { HelixStreamData } from '@twurple/api/lib/index.js';
import axios from 'axios';
import { mk } from 'date-fns/locale';

import { app } from '../../index.js';
import { DefaultCommand } from '../types.js';

const prisma = app.get(PrismaService);
const staticApi = app.get(TwitchApiService);
const redis = app.get(RedisService);

const messages = Object.freeze({
  GAME_NOT_FOUND: 'Game not found.',
  NO_ACCOUNTS: 'You have not added account.',
});


const dotaApi = axios.create({
  baseURL: `http://api.steampowered.com/`,
  timeout: 1000,
});

dotaApi.interceptors.request.use((req) => {
  req.params = req.params || {};
  req.params['key'] = config.STEAM_API_KEY;
  return req;
});

const getGames = async (accounts: string[]) => {
  const rps = await Promise.all(accounts.map(a => redis.get(`dotaRps:${a}`)));
  if (!rps.filter(r => r !== null).length) {
    return messages.GAME_NOT_FOUND;
  }
  const cachedRps = rps.filter(r => r !== null).map(r => JSON.parse(r!));

  const cachedGames = await Promise.all(accounts.map(a => redis.get(`dotaMatches:${a}`)));
  if (!cachedGames.filter(r => r !== null).length) {
    return messages.GAME_NOT_FOUND;
  }

  const parsedGames = cachedGames.filter(g => g !== null).map(g => JSON.parse(g!) as DotaGame);
  const dbGames = await prisma.dotaMatch.findMany({
    where: {
      lobbyId: {
        in: cachedRps.map(r => r.lobbyId),
      },
      players: {
        hasSome: accounts.map(a => Number(a)),
      },
    },
    orderBy: {
      startedAt: 'desc',
    },
    include: {
      gameMode: true,
    },
    take: 2,
  });

  if (!dbGames.length) return messages.GAME_NOT_FOUND;

  return dbGames.map(g => {
    const cachedGame = parsedGames.find(game => game.match_id === g.match_id)!;

    return {
      ...g,
      players: cachedGame.players,
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
    name: 'np',
    permission: 'VIEWER',
    visible: true,
    handler: async (state) => {
      if (!state.channelId) return;

      const accounts = await getAccounts(state.channelId);
      if (typeof accounts === 'string') return accounts;

      const games = await getGames(accounts.map(a => a.id));
      if (typeof games === 'string') return games;

      return games
        .map(g => {
          const avgMmr = g.gameMode.id === 22 ? ` (${g.avarage_mmr}mmr)` : '';
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
      if (!stream) return 'Stream is offline';
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
          finished: true,
        },
      }).then(gms => gms.filter(g => !g.players_heroes.some(h => h === 0)));


      const matchesRequest = await Promise.all(games.map(g => dotaApi.get('IDOTA2Match_570/GetMatchDetails/v1', { params: { match_id: g.match_id } })))
        .then(requests => requests.filter(r => r.status === 200));

      const matchesData: any[] = [];
      for (const response of matchesRequest) {
        if (response.data) matchesData.push(response.data);
      }

      const matchesByGameMode: { [x: number]: any[] } = {};
      gameModes.forEach(m => {
        matchesByGameMode[m.id] = [];
      });

      for (const account of accounts) {
        for (const match of matchesData) {
          if (typeof match.result.radiant_win === 'undefined' || !match.result?.players) continue;

          await prisma.dotaMatch.update({
            where: {
              match_id: match.result.match_id.toString(),
            },
            data: {
              finished: true,
            },
          });
          const player = match.result.players.find((p: any) => p.account_id === Number(account.id));
          if (!player) continue;
          const isWinner = player.team_number === 0 && match.result.radiant_win;
          const hero = dotaHeroes.find(h => h.id === player.hero_id);
          matchesByGameMode[match.result.game_mode]?.push({
            isWinner,
            hero,
            kills: player.kills,
            deaths: player.deaths,
            assists: player.assists,
          });
        }
      }

      const result: string[] = [];

      for (const [modeId, matches] of Object.entries(matchesByGameMode).filter(e => e[1].length)) {
        const wins = matches.filter(r => r.isWinner);
        const mode = gameModes.find(m => m.id === Number(modeId));
        const heroesResult = matches.map(m => `${m.hero.localized_name}(${m.isWinner ? 'W' : 'L'}) [${m.kills}/${m.deaths}/${m.assists}]`);
        let msg = `${mode?.name ?? 'Unknown'} W ${wins.length} — L ${matches.length - wins.length}`;
        if (mode?.id === 22) msg += `: ${heroesResult.join(', ')} `;
        result.push(msg);
      }

      return result.length ? result.join(' | ') : 'W 0 — L 0';
    },
  },
];