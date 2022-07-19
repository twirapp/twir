import { Module } from '@nestjs/common';
import { ClientsModule, Transport } from '@nestjs/microservices';
import { config } from '@tsuwari/config';


@Module({
  providers: [
    {
      provide: 'NATS',
      useValue: ClientsModule.register([
        { name: 'NATS', transport: Transport.NATS, options: { servers: [config.NATS_URL] } },
      ]),
    },
  ],
  exports: ['NATS'],
})
export class NatsModule { }
