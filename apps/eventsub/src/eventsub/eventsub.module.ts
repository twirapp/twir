import { Module } from '@nestjs/common';
import { config } from '@tsuwari/config';
import { TwitchApiService } from '@tsuwari/shared';
import * as ngrok from 'ngrok';

import { HandlerModule } from '../handler/handler.module.js';
import { EventSubController } from './eventsub.controller.js';
import { EventSub } from './eventsub.service.js';

@Module({
  controllers: [EventSubController],
  imports: [
    HandlerModule,
  ],
  providers: [TwitchApiService],
})
export class EventSubModule {
  static async register() {
    let hostName = '';

    if (config.isDev) {
      const tunnel = await ngrok.connect(3003);
      hostName = tunnel.replace('https://', '');
    } else {
      hostName = `eventsub.${config.HOSTNAME.replace('https://', '')}`;
    }

    return {
      module: EventSubModule,
      providers: [
        {
          provide: 'HOSTNAME',
          useValue: hostName,
        },
        EventSub,
      ],
      exports: [EventSub],
    };
  }
}