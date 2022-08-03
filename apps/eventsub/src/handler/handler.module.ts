import { Module } from '@nestjs/common';
import { ClientsModule, Transport } from '@nestjs/microservices';
import { config } from '@tsuwari/config';
import { TwitchApiService } from '@tsuwari/shared';

import { HandlerService } from './handler.service.js';

@Module({
  imports: [
    ClientsModule.register([
      { name: 'NATS', transport: Transport.NATS, options: { servers: [config.NATS_URL] } },
    ]),
  ],
  providers: [TwitchApiService, HandlerService],
  exports: [HandlerService],
})
export class HandlerModule { }
