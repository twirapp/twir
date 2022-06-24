import { Controller, Logger } from '@nestjs/common';
import { EventPattern, MessagePattern, Payload } from '@nestjs/microservices';
import { ClientProxyResult, ClientProxyCommandsKey, ClientProxyEventsKey, ClientProxyEvents } from '@tsuwari/shared';
import { of } from 'rxjs';

import { AppService } from './app.service.js';

@Controller()
export class AppController {
  private readonly logger = new Logger('StreamStatus');

  constructor(private readonly appService: AppService) { }

  @MessagePattern<ClientProxyCommandsKey>('streamstatuses.process')
  async cacheStreams(data: string[]): Promise<ClientProxyResult<'streamstatuses.process'>> {
    this.logger.log(`Starting to process ${data.length} streams`);
    return of(await this.appService.handleChannels(data));
  }

  @EventPattern<ClientProxyEventsKey>('streams.online')
  streamOnline(@Payload() data: ClientProxyEvents['streams.online']['input']) {
    this.logger.log(`Starting to process ${data.channelId} online`);
    this.appService.handleOnline(data);
  }

  @EventPattern<ClientProxyEventsKey>('streams.offline')
  streamOffline(@Payload() data: ClientProxyEvents['streams.offline']['input']) {
    this.logger.log(`Starting to process ${data.channelId} offline`);
    this.appService.handleOffline(data);
  }

  @EventPattern<ClientProxyEventsKey>('stream.update')
  streamUpdate(@Payload() data: ClientProxyEvents['stream.update']['input']) {
    this.logger.log(`Starting to process ${data.broadcaster_user_id} update`);
    this.appService.handleUpdate(data);
  }
}
