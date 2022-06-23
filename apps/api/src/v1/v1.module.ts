import { Module } from '@nestjs/common';

import { CommandsModule } from './commands/commands.module.js';
import { GreetingsModule } from './greetings/greetings.module.js';
import { FaceitModule } from './integrations/faceit/faceit.module.js';
import { LastfmModule } from './integrations/lastfm/lastfm.module.js';
import { SpotifyModule } from './integrations/spotify/spotify.module.js';
import { VkModule } from './integrations/vk/vk.module.js';
import { KeywordsModule } from './keywords/keywords.module.js';
import { ModerationModule } from './moderation/moderation.module.js';
import { SettingsModule } from './settings/settings.module.js';
import { StreamsModule } from './streams/streams.module.js';
import { TimersModule } from './timers/timers.module.js';
import { VariablesModule } from './variables/variables.module.js';

@Module({
  imports: [CommandsModule, GreetingsModule, TimersModule, SpotifyModule, LastfmModule, KeywordsModule, VariablesModule, ModerationModule, StreamsModule, SettingsModule, VkModule, FaceitModule],
})
export class V1Module { }
