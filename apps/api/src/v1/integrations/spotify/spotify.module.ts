import { Module } from '@nestjs/common';

import { SpotifyIntegrationService } from './integration.js';
import { SpotifyController } from './spotify.controller.js';
import { SpotifyService } from './spotify.service.js';

@Module({
  controllers: [SpotifyController],
  providers: [SpotifyIntegrationService, SpotifyService],
})
export class SpotifyModule { }
