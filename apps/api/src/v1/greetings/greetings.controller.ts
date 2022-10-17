import {
  Body,
  CacheTTL,
  Controller,
  Delete,
  Get,
  Param,
  Post,
  Put,
  UseGuards,
  UseInterceptors,
} from '@nestjs/common';
import Express from 'express';

import { DashboardAccessGuard } from '../../guards/DashboardAccess.guard.js';
import { CustomCacheInterceptor } from '../../helpers/customCacheInterceptor.js';
import { JwtAuthGuard } from '../../jwt/jwt.guard.js';
import { GreetingCreateDto } from './dto/create.js';
import { GreetingsService } from './greetings.service.js';

@Controller('v1/channels/:channelId/greetings')
export class GreetingsController {
  constructor(private readonly greetingsService: GreetingsService) {}

  @UseGuards(JwtAuthGuard, DashboardAccessGuard)
  @Get()
  root(@Param('channelId') channelId: string) {
    return this.greetingsService.getList(channelId);
  }

  @UseGuards(JwtAuthGuard, DashboardAccessGuard)
  @Post()
  async create(@Param('channelId') channelId: string, @Body() body: GreetingCreateDto) {
    const result = await this.greetingsService.create(channelId, body);
    return result;
  }

  @UseGuards(JwtAuthGuard, DashboardAccessGuard)
  @Put(':greetingId')
  async update(
    @Param('channelId') channelId: string,
    @Param('greetingId') greetingId: string,
    @Body() body: GreetingCreateDto,
  ) {
    const result = await this.greetingsService.update(channelId, greetingId, body);
    return result;
  }

  @UseGuards(JwtAuthGuard, DashboardAccessGuard)
  @Delete(':greetingId')
  async delete(@Param('channelId') channelId: string, @Param('greetingId') greetingId: string) {
    const result = await this.greetingsService.delete(channelId, greetingId);
    return result;
  }
}
