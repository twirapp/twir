import { Module } from '@nestjs/common';
import { ClientsModule, Transport } from '@nestjs/microservices';
import { config } from '@tsuwari/config';

import { StreamStatusService } from './statuschecker.service.js';

@Module({
  imports: [
    ClientsModule.register([
      { name: 'NATS', transport: Transport.NATS, options: { servers: [config.NATS_URL] } },
    ]),
  ],
  providers: [StreamStatusService],
})
export class StreamStatusModule { }
