import { Global, Module } from '@nestjs/common';
import { ClientProxyFactory, Transport } from '@nestjs/microservices';
import { config } from '@tsuwari/config';
import { resolveProtoPath } from '@tsuwari/grpc';

@Global()
@Module({
  providers: [
    {
      provide: 'STREAMSTATUS_MICROSERVICE',
      useValue: ClientProxyFactory
        .create({
          transport: Transport.GRPC,
          options: {
            package: 'StreamStatus',
            protoPath: resolveProtoPath('streamstatus'),
            url: config.MICROSERVICE_STREAM_STATUS_URL,
          },
        }),
    },
    {
      provide: 'WATCHED_MICROSERVICE',
      useValue: ClientProxyFactory
        .create({
          transport: Transport.GRPC,
          options: {
            package: 'Watched',
            protoPath: resolveProtoPath('watched'),
            url: config.MICROSERVICE_WATCHED_URL,
          },
        }),
    },
  ],
  exports: ['STREAMSTATUS_MICROSERVICE', 'WATCHED_MICROSERVICE'],
})
export class MicroservicesModule { }
