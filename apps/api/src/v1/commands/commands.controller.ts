import { Body, CacheTTL, CACHE_MANAGER, Controller, Delete, Get, Inject, Param, Post, Put, UseGuards, UseInterceptors } from '@nestjs/common';
import { Cache } from 'cache-manager';
import { Request } from 'express';

import { DashboardAccessGuard } from '../../guards/DashboardAccess.guard.js';
import { CustomCacheInterceptor } from '../../helpers/customCacheInterceptor.js';
import { JwtAuthGuard } from '../../jwt/jwt.guard.js';
import { CommandsService } from './commands.service.js';
import { UpdateOrCreateCommandDto } from './dto/create.js';

@Controller('v1/channels/:channelId/commands')
export class CommandsController {
  constructor(
    private readonly commandsSerivce: CommandsService,
    @Inject(CACHE_MANAGER) private readonly cacheManager: Cache,
  ) { }

  private delCache(channelId: string) {
    this.cacheManager.del(`nest:cache:v1/channels/${channelId}/commands`);
  }

  @CacheTTL(600)
  @UseInterceptors(CustomCacheInterceptor(ctx => {
    const req = ctx.switchToHttp().getRequest() as Request;
    return `nest:cache:v1/channels/${req.params.channelId}/commands`;
  }))
  @UseGuards(JwtAuthGuard, DashboardAccessGuard)
  @Get()
  root(@Param('channelId') channelId: string) {
    return this.commandsSerivce.getList(channelId);
  }

  @UseGuards(JwtAuthGuard, DashboardAccessGuard)
  @Post()
  async create(@Param('channelId') channelId: string, @Body() body: UpdateOrCreateCommandDto) {
    const result = await this.commandsSerivce.create(channelId, body);
    this.delCache(channelId);
    return result;
  }

  @UseGuards(JwtAuthGuard, DashboardAccessGuard)
  @Delete(':commandId')
  async delete(@Param('channelId') channelId: string, @Param('commandId') commandId: string) {
    const result = await this.commandsSerivce.delete(channelId, commandId);
    this.delCache(channelId);
    return result;
  }

  @UseGuards(JwtAuthGuard, DashboardAccessGuard)
  @Put(':commandId')
  async update(@Param('channelId') channelId: string, @Param('commandId') commandId: string, @Body() body: UpdateOrCreateCommandDto) {
    const result = this.commandsSerivce.update(channelId, commandId, body);
    this.delCache(channelId);
    return result;
  }
}
