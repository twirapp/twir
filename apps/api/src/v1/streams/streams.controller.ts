import { Controller, Get, Param, UseGuards } from '@nestjs/common';
import { Payload } from '@nestjs/microservices';
import { type ClientProxyEventPayload, EventPattern } from '@tsuwari/shared';

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
  handler(@Payload() data: ClientProxyEventPayload<'streams.online'> | ClientProxyEventPayload<'streams.offline'>) {
    this.streamsService.handleStreamStateChange(data.channelId);
  }
}
