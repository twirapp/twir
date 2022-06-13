import { Global, Module } from '@nestjs/common';
import { ClientProxyFactory, Transport } from '@nestjs/microservices';
import { config } from '@tsuwari/config';

@Global()
@Module({
  providers: [],
  exports: [],
})
export class MicroservicesModule { }
