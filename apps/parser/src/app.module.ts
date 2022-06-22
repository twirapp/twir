import { Module } from '@nestjs/common';
import { PrismaModule } from '@tsuwari/prisma';
import { RedisService, TwitchApiService } from '@tsuwari/shared';


import { AppController } from './app.controller.js';
import { AppService } from './app.service.js';
import { HelpersService } from './helpers.service.js';
import { FaceitIntegration } from './integrations/faceit.js';
import { LastFmIntegration } from './integrations/lastfm.js';
import { SpotifyIntegration } from './integrations/spotify.js';
import { VkIntegration } from './integrations/vk.js';
import { ParserCache } from './variables/cache.js';
import { VariablesParser } from './variables/index.js';

@Module({
  imports: [PrismaModule],
  controllers: [AppController],
  providers: [RedisService, ParserCache, VariablesParser, LastFmIntegration, SpotifyIntegration, VkIntegration, FaceitIntegration, TwitchApiService, HelpersService, AppService],
})
export class AppModule { }
