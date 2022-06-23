import { Controller, Get, Param, UseGuards } from '@nestjs/common';
import { EventPattern, Payload } from '@nestjs/microservices';

import { DashboardAccessGuard } from '../../guards/DashboardAccess.guard.js';
import { JwtAuthGuard } from '../../jwt/jwt.guard.js';
import { StreamsService } from './streams.service.js';

@Controller('v1/channels/:channelId/streams')
export class StreamsController {
  constructor(private readonly streamsService: StreamsService) { }

  @Get()
  @UseGuards(JwtAuthGuard, DashboardAccessGuard)
  getStream(@Param('channelId') channelId: string) {
    return this.streamsService.getStream(channelId);
  }

  @EventPattern('streams.online')
  @EventPattern('streams.offline')
  handler(@Payload() data: { channelId: string, streamId?: string }) {
    this.streamsService.handleStreamStateChange(data.channelId);
  }
}
