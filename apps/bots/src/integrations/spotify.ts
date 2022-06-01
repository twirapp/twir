import { SpotifyIntegration as Spotify } from '@tsuwari/spotify';

import { prisma } from '../libs/prisma.js';

export const SpotifyIntegration = new Spotify(prisma);
