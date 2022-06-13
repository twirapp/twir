import { Controller, Logger } from '@nestjs/common';
import { MessagePattern } from '@nestjs/microservices';

import { AppService } from './app.service.js';

@Controller()
export class AppController {
  private readonly logger = new Logger('StreamStatus');

  constructor(private readonly appService: AppService) { }

  @MessagePattern('streamstatuses.process')
  async cacheStreams(data: string[]) {
    this.logger.log(`Starting to process ${data.length} streams`);
    return await this.appService.handleChannels(data);
  }
}
