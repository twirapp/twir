import { UserIntegration } from '@tsuwari/spotify';

import { LastFMDbData, LastFmIntegration } from '../../integrations/lastfm.js';
import { SpotifyIntegration } from '../../integrations/spotify.js';
import { VKDBData, VKIntegration } from '../../integrations/vk.js';
import { prisma } from '../../libs/prisma.js';
import { Module } from '../index.js';

export const song: Module = {
  key: 'currentsong',
  handler: async (_, state) => {
    const enabledIntegrations = await prisma.channelIntegration.findMany({
      where: {
        channelId: state.channelId, 
        enabled: true,
      },
      include: {
        integration: true,
      },
    });

    let result = '';

    for (const integration of enabledIntegrations) {
      switch (integration.integration.service) {
        case 'SPOTIFY': {
          if (!integration.accessToken || !integration.refreshToken) continue;
          const instance = await new UserIntegration(state.channelId, SpotifyIntegration, prisma).init();
          if (!instance) continue;
          const song = await instance.getCurrentSong();
          if (song) result = song;
          
          break;
        }
        case 'VK': {
          const data = integration.data as unknown as VKDBData;
          if (!data.userId) continue;
          const song = await VKIntegration.fetchSong(data.userId);
          if (song) result = song;
          
          break;
        }
        case 'LASTFM': {
          const data = integration.data as unknown as LastFMDbData;
          if (!data.username) continue;
          const song = await LastFmIntegration.fetchSong(data.username);
          if (song) {
            result = song;
          }
          break;
        }

        default:
          break;
      }
    }

    if (result && result.length) {
      return result;
    } else {
      return 'Not playing any song.';
    }
  },
};
