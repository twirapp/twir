import {
  Body,
  CACHE_MANAGER,
  Controller,
  Delete,
  Get,
  Inject,
  Param,
  Post,
  Put,
  UseGuards,
} from '@nestjs/common';
import CacheManager from 'cache-manager';

import { DashboardAccessGuard } from '../../guards/DashboardAccess.guard.js';
import { JwtAuthGuard } from '../../jwt/jwt.guard.js';
import { CommandsService } from './commands.service.js';
import { UpdateOrCreateCommandDto } from './dto/create.js';

@Controller('v1/channels/:channelId/commands')
export class CommandsController {
  constructor(private readonly commandsSerivce: CommandsService) {}

  @UseGuards(JwtAuthGuard, DashboardAccessGuard)
  @Get()
  root(@Param('channelId') channelId: string) {
    return this.commandsSerivce.getList(channelId);
  }

  @UseGuards(JwtAuthGuard, DashboardAccessGuard)
  @Post()
  async create(@Param('channelId') channelId: string, @Body() body: UpdateOrCreateCommandDto) {
    const result = await this.commandsSerivce.create(channelId, body);
    return result;
  }

  @UseGuards(JwtAuthGuard, DashboardAccessGuard)
  @Delete(':commandId')
  async delete(@Param('channelId') channelId: string, @Param('commandId') commandId: string) {
    const result = await this.commandsSerivce.delete(channelId, commandId);
    return result;
  }

  @UseGuards(JwtAuthGuard, DashboardAccessGuard)
  @Put(':commandId')
  async update(
    @Param('channelId') channelId: string,
    @Param('commandId') commandId: string,
    @Body() body: UpdateOrCreateCommandDto,
  ) {
    const result = this.commandsSerivce.update(channelId, commandId, body);
    return result;
  }
}
