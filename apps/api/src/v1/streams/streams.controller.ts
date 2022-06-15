import { Controller } from '@nestjs/common';
import { EventPattern, Payload } from '@nestjs/microservices';

import { StreamsService } from './streams.service.js';

@Controller()
export class StreamsController {
  constructor(private readonly streamsService: StreamsService) { }

  @EventPattern('streams.online')
  @EventPattern('streams.offline')
  handler(@Payload() data: { channelId: string, streamId?: string }) {
    this.streamsService.handleStreamStateChange(data.channelId);
  }
}
