import {
  Body,
  CacheTTL,
  Controller,
  Get,
  Param,
  Patch,
  UseGuards,
  UseInterceptors,
} from '@nestjs/common';
import Express from 'express';

import { DashboardAccessGuard } from '../../guards/DashboardAccess.guard.js';
import { CustomCacheInterceptor } from '../../helpers/customCacheInterceptor.js';
import { JwtAuthGuard } from '../../jwt/jwt.guard.js';
import { BotService } from './bot.service.js';
import { JoinOrLeaveDto } from './dto/joinOrLeave.js';

@Controller('v1/channels/:channelId/bot')
export class BotController {
  constructor(private readonly service: BotService) {}

  @CacheTTL(10)
  @UseInterceptors(
    CustomCacheInterceptor((ctx) => {
      const req = ctx.switchToHttp().getRequest() as Express.Request;
      return `nest:cache:socket/isBotMod:${req.params.channelId}`;
    }),
  )
  @UseGuards(JwtAuthGuard, DashboardAccessGuard)
  @Get('checkmod')
  mod(@Param('channelId') channelId: string) {
    return this.service.isBotMod(channelId);
  }

  @UseGuards(JwtAuthGuard, DashboardAccessGuard)
  @Patch('connection')
  joinOrLeave(@Param('channelId') channelId: string, @Body() body: JoinOrLeaveDto) {
    return this.service.botJoinOrLeave(body.action, channelId);
  }
}
