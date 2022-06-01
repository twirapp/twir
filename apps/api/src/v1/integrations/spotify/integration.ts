
import { Injectable } from '@nestjs/common';
import { PrismaService } from '@tsuwari/prisma';
import { SpotifyIntegration } from '@tsuwari/spotify';


@Injectable()
export class SpotifyIntegrationService extends SpotifyIntegration {
  constructor(prisma: PrismaService) {
    super(prisma);
  }
}
