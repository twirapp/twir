import { Injectable } from '@nestjs/common';
import { Client, Transport } from '@nestjs/microservices';
import { config } from '@tsuwari/config';
import { ClientProxy, ClientProxyCommands } from '@tsuwari/shared';

@Injectable()
export class ParserService {
  @Client({ transport: Transport.NATS, options: { servers: [config.NATS_URL] } })
  nats: ClientProxy;

  parseChatMessage(raw: string) {
    return this.nats.send('parseChatMessage', raw).toPromise();
  }

  parseResponse(input: ClientProxyCommands['parseResponse']['input']) {
    return this.nats.send('parseResponse', input).toPromise();
  }
}
