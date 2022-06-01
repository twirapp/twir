import { Global, Module } from '@nestjs/common';
import { ClientProxyFactory, Transport } from '@nestjs/microservices';
import { config } from '@tsuwari/config';
import { resolveProtoPath } from '@tsuwari/grpc';

@Global()
@Module({
  providers: [
    {
      provide: 'BOTS_MICROSERVICE',
      useValue: ClientProxyFactory
        .create({
          transport: Transport.GRPC,
          options: {
            package: 'Bots',
            protoPath: resolveProtoPath('bots'),
            url: config.MICROSERVICE_BOTS_URL,
          },
        }),
    },
  ],
  exports: ['BOTS_MICROSERVICE'],
})
export class BotsMicroserviceModule { }
