import { endOfDay, startOfDay } from 'date-fns';
import omit from 'lodash.omit';

import { prisma } from '../libs/prisma.js';

export interface FaceitDbData {
  username: string;
  game?: string;
}

class FaceitIntegrationClass {
  #apiKey: string;
  readonly #baseApiUrl = 'https://open.faceit.com/data/v4/';

  async init() {
    const service = await prisma.integration.findFirst({
      where: {
        service: 'FACEIT',
      },
    });

    if (service && service.apiKey) {
      this.#apiKey = service.apiKey;
    }

    return this;
  }

  async fetchStats(nickname: string, game = 'csgo') {
    const baseProfileRequest = await fetch(`${this.#baseApiUrl}/players?nickname=${nickname}`);
    if (!baseProfileRequest.ok) return null;

    const baseProfile = await baseProfileRequest.json() as unknown as BaseProfileResponse;
    const userGame = baseProfile.games[game];
    if (!userGame) return null;

    const { faceit_elo: elo, skill_level: lvl } = userGame;

    let [latestMatches, userStats] = await Promise.all([
      this.#getLastMatches(baseProfile.player_id, game),
      this.#getUserStats(baseProfile.player_id, game),
    ]);

    if (!latestMatches) {
      latestMatches = {
        todayEloDiff: 'cannot fetch data',
        latestMatches: [],
      };
    }

    return {
      elo,
      lvl,
      todayEloDiff: latestMatches.todayEloDiff,
      latestMatches: latestMatches.latestMatches,
      latestMatchesTrend: {
        simple: latestMatches.latestMatches
          .map((m) => {
            return m.result === 'Won' ? 'W' : 'L';
          })
          .join(''),
        extended: latestMatches.latestMatches
          .map((m) => {
            return `${m.result === 'Won' ? 'W' : 'L'} ${m.eloDiff}`;
          })
          .join(' | '),
        score: {
          loses: latestMatches.latestMatches.filter((m) => m.result !== 'Won').length,
          wins: latestMatches.latestMatches.filter((m) => m.result === 'Won').length,
        },
      },
      stats: {
        lifetime: userStats?.lifetime,
      },
    };
  }

  async #getLastMatches(playerId: string, game = 'csgo') {
    const matchesRequest = await fetch(
      `https://api.faceit.com/stats/api/v1/stats/time/users/${playerId}/games/${game}?size=30`,
    );

    if (!matchesRequest.ok) return null;

    const matchesResponse = (await matchesRequest.json()) as unknown as Array<FaceitMatch>;

    const dayStart = startOfDay(Date.now()).getTime();
    const dayEnd = endOfDay(Date.now()).getTime();

    const matches = matchesResponse.map((m) => ({
      ...m,
      elo: Number(m.elo),
      result: m.i10 === '1' ? 'Won' : 'Lost',
    }));

    const result = [];

    for (const match of matches) {
      if (match.date > dayStart && match.date < dayEnd) {
        const index = matches.indexOf(match);
        const prev = matches[index + 1];

        if (prev) {
          result.push(prev.elo ? (prev.elo > match.elo ? -(prev.elo - match.elo) : match.elo - prev.elo) : 0);
        }
      }
    }

    const todayEloDiff = result.reduce((prev, curr) => prev + curr, 0);

    return {
      todayEloDiff: todayEloDiff > 0 ? `+${todayEloDiff}` : `${todayEloDiff}`,
      latestMatches: matches
        .filter((m) => new Date(m.created_at).getTime() > dayStart)
        .map((m) => {
          const prev = matches[matches.indexOf(m) + 1];

          return {
            team: m.i5,
            teamScore: m.result === 'Won' ? m.i18.split(' / ').reverse().join(' / ') : m.i18,
            map: m.i1,
            kd: m.c2,
            hs: {
              percentage: m.c4,
              number: m.i13,
            },
            eloDiff: prev?.elo ? (prev.elo > m.elo ? `-${prev.elo - m.elo}` : `+${m.elo - prev.elo}`) : '0',
            kills: m.i6,
            death: m.i8,
            result: m.result as 'Won' | 'Lose',
            createdAt: new Date(m.created_at),
            updatedAt: new Date(m.updated_at),
          };
        }),
    };
  }

  async #getUserStats(playerId: string, game = 'csgo') {
    const statsRequest = await fetch(`${this.#baseApiUrl}players/${playerId}/stats/${game}`);

    if (!statsRequest.ok) return null;
    const data = await statsRequest.json() as FaceitUserStats;
    return omit(data, 'lifetime.Recent Results');
  }
}

export const FaceitIntegration = await new FaceitIntegrationClass().init();

type BaseProfileResponse = {
  games: {
    [x: string]: {
      skill_level: number;
      faceit_elo: number;
    };
  };
  player_id: string;
};

type FaceitMatch = {
  elo: string;
  date: number;
  i18: string;
  i5: string;
  i1: string;
  c2: string;
  c4: string;
  c5: string;
  i13: string;
  i6: string;
  i8: string;
  i10: string;
  created_at: number;
  updated_at: number;
};

export type FaceitUserStats = {
  lifetime: {
    Matches: string;
    'Longest Win Streak': string;
    'Win Rate %': string;
    'Total Headshots %': string;
    'Average K/D Ratio': string;
    'Average Headshots %': string;
    'K/D Ratio': string;
    'Current Win Streak': string;
    Wins: string;
  };
};
