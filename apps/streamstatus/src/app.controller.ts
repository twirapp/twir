import { Controller, Logger } from '@nestjs/common';
import { MessagePattern } from '@nestjs/microservices';
import { ClientProxyResult } from '@tsuwari/shared';
import { of } from 'rxjs';

import { AppService } from './app.service.js';

@Controller()
export class AppController {
  private readonly logger = new Logger('StreamStatus');

  constructor(private readonly appService: AppService) { }

  @MessagePattern('streamstatuses.process')
  async cacheStreams(data: string[]): Promise<ClientProxyResult<'streamstatuses.process'>> {
    this.logger.log(`Starting to process ${data.length} streams`);
    return of(await this.appService.handleChannels(data));
  }
}
