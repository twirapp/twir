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
  ],
  exports: ['STREAMSTATUS_MICROSERVICE'],
})
export class MicroservicesModule { }
