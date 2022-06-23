import { Injectable } from '@nestjs/common';
import { PrismaService } from '@tsuwari/prisma';
import { SpotifyIntegration as Spotify } from '@tsuwari/spotify';


@Injectable()
export class SpotifyIntegration extends Spotify {
  constructor(prisma: PrismaService) {
    super(prisma);
  }
}

