import { Body, CacheTTL, CACHE_MANAGER, Controller, Delete, Get, Inject, Param, Post, Put, UseGuards, UseInterceptors } from '@nestjs/common';
import { Cache } from 'cache-manager';
import { Request } from 'express';

import { DashboardAccessGuard } from '../../guards/DashboardAccess.guard.js';
import { CustomCacheInterceptor } from '../../helpers/customCacheInterceptor.js';
import { JwtAuthGuard } from '../../jwt/jwt.guard.js';
import { CreateVariableDto } from './dto/create.js';
import { VariablesService } from './variables.service.js';

@Controller('v1/channels/:channelId/variables')
export class VariablesController {
  constructor(
    private readonly variablesService: VariablesService,
    @Inject(CACHE_MANAGER) private cacheManager: Cache,
  ) { }

  private async delCache(channelId: string) {
    await this.cacheManager.del(`nest:cache:v1/channels/${channelId}/variables`);
  }

  @Get('builtin')
  @UseGuards(JwtAuthGuard, DashboardAccessGuard)
  builtIn() {
    return this.variablesService.getBuildInVariables();
  }

  @CacheTTL(600)
  @UseInterceptors(CustomCacheInterceptor(ctx => {
    const req = ctx.switchToHttp().getRequest() as Request;
    return `nest:cache:v1/channels/${req.params.channelId}/variables`;
  }))
  @UseGuards(JwtAuthGuard, DashboardAccessGuard)
  @Get()
  root(@Param('channelId') channelId: string) {
    return this.variablesService.getList(channelId);
  }

  @UseGuards(JwtAuthGuard, DashboardAccessGuard)
  @Post()
  async create(@Param('channelId') channelId: string, @Body() body: CreateVariableDto) {
    const result = await this.variablesService.create(channelId, body);
    await this.delCache(channelId);
    return result;
  }

  @UseGuards(JwtAuthGuard, DashboardAccessGuard)
  @Delete(':id')
  async delete(@Param('channelId') channelId: string, @Param('id') id: string) {
    const result = await this.variablesService.delete(channelId, id);
    await this.delCache(channelId);
    return result;
  }

  @UseGuards(JwtAuthGuard, DashboardAccessGuard)
  @Put(':id')
  async update(@Param('channelId') channelId: string, @Param('id') id: string, @Body() body: CreateVariableDto) {
    const result = await this.variablesService.update(channelId, id, body);
    await this.delCache(channelId);
    return result;
  }
}
