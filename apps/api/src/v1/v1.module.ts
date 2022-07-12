import { Module } from '@nestjs/common';
import { config } from '@tsuwari/config';
import { S3Module } from 'nestjs-s3';

import { CommandsModule } from './commands/commands.module.js';
import { FeedbackModule } from './feedback/feedback.module.js';
import { FilesModule } from './files/files.module.js';
import { GreetingsModule } from './greetings/greetings.module.js';
import { FaceitModule } from './integrations/faceit/faceit.module.js';
import { LastfmModule } from './integrations/lastfm/lastfm.module.js';
import { SpotifyModule } from './integrations/spotify/spotify.module.js';
import { VkModule } from './integrations/vk/vk.module.js';
import { KeywordsModule } from './keywords/keywords.module.js';
import { ModerationModule } from './moderation/moderation.module.js';
import { NotificationsModule } from './notifications/notifications.module.js';
import { SettingsModule } from './settings/settings.module.js';
import { StreamsModule } from './streams/streams.module.js';
import { TimersModule } from './timers/timers.module.js';
import { VariablesModule } from './variables/variables.module.js';

@Module({
  imports: [
    S3Module.forRoot({
      config: {
        accessKeyId: config.MINIO_ACCESS_KEY,
        secretAccessKey: config.MINIO_ACCESS_SECRET,
        endpoint: config.MINIO_URL,
        s3ForcePathStyle: true,
        region: 'eu-east-1',
        signatureVersion: 'v4',
      },
    }), CommandsModule, GreetingsModule, TimersModule, SpotifyModule, LastfmModule, KeywordsModule, VariablesModule, ModerationModule, StreamsModule, SettingsModule, VkModule, FaceitModule, NotificationsModule, FeedbackModule, FilesModule],
})
export class V1Module { }
