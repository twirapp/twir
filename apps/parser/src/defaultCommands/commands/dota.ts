import { PrismaService } from '@tsuwari/prisma';
import { DotaGame, RedisService, TwitchApiService } from '@tsuwari/shared';

import { app } from '../../index.js';
import { DefaultCommand } from '../types.js';

const prisma = app.get(PrismaService);
const staticApi = app.get(TwitchApiService);
const redis = app.get(RedisService);

const messages = Object.freeze({
  GAME_NOT_FOUND: 'Game not found.'
})

const getGames = async (accounts: string[]) => {
  const rps = await Promise.all(accounts.map(a => redis.get(`dotaRps:${a}`)))
  if (!rps.filter(r => r !== null).length) {
    return messages.GAME_NOT_FOUND;
  }

  const cachedGames = await Promise.all(accounts.map(a => redis.get(`dotaMatches:${a}`)))
  if (!cachedGames.filter(r => r !== null).length) {
    return messages.GAME_NOT_FOUND;
  }

  const parsedGames = cachedGames.map(g => JSON.parse(g!) as DotaGame)
  const dbGames = await prisma.dotaMatch.findMany({
    where: {
      startedAt: {
        gte: new Date(new Date().getTime() - 900000)
      },
      players: {
        hasSome: accounts.map(a => Number(a))
      }
    },
    orderBy: {
      startedAt: 'desc'
    },
    take: 2,
  })

  return dbGames.map(g => {
    const cachedGame = parsedGames.find(game => game.match_id === g.match_id)!

    return {
      ...g,
      players: cachedGame.players
    }
  })
}

const dota: DefaultCommand[] = [{
  name: 'np',
  permission: 'VIEWER',
  visible: false,
  handler: async (state, params) => {
    if (!state.channelId) return;

    const accounts = await prisma.dotaAccount.findMany({
      where: {
        channelId: state.channelId,
      }
    })

    if (!accounts.length) return 'You have not added account.'

    const games = await getGames(accounts.map(a => a.id))
  },
}];