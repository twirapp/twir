import { Controller, Logger } from '@nestjs/common';
import { Payload } from '@nestjs/microservices';
import { ClientProxyResult, MessagePattern, type ClientProxyCommandPayload, EventPattern, type ClientProxyEventPayload } from '@tsuwari/shared';
import { of } from 'rxjs';

import { AppService } from './app.service.js';

@Controller()
export class AppController {
  private readonly logger = new Logger('StreamStatus');

  constructor(private readonly appService: AppService) { }

  @MessagePattern('streamstatuses.process')
  async cacheStreams(@Payload() data: ClientProxyCommandPayload<'streamstatuses.process'>): Promise<ClientProxyResult<'streamstatuses.process'>> {
    this.logger.log(`Starting to process ${data.length} streams`);
    return of(await this.appService.handleChannels(data));
  }

  @EventPattern('streams.online')
  streamOnline(@Payload() data: ClientProxyEventPayload<'streams.online'>) {
    this.logger.log(`Starting to process ${data.channelId} online`);
    this.appService.handleOnline(data);
  }

  @EventPattern('streams.offline')
  streamOffline(@Payload() data: ClientProxyEventPayload<'streams.offline'>) {
    this.logger.log(`Starting to process ${data.channelId} offline`);
    this.appService.handleOffline(data);
  }

  @EventPattern('stream.update')
  streamUpdate(@Payload() data: ClientProxyEventPayload<'stream.update'>) {
    this.logger.log(`Starting to process ${data.broadcaster_user_id} update`);
    this.appService.handleUpdate(data);
  }
}
