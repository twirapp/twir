import { Body, Controller, Delete, Get, Param, Post, Put, UseGuards } from '@nestjs/common';

import { DashboardAccessGuard } from '../../guards/DashboardAccess.guard.js';
import { JwtAuthGuard } from '../../jwt/jwt.guard.js';
import { CreateTimerDto } from './dto/create.js';
import { TimersService } from './timers.service.js';

@Controller('v1/channels/:channelId/timers')
export class TimersController {
  constructor(private readonly timersService: TimersService) {}

  @UseGuards(JwtAuthGuard, DashboardAccessGuard)
  @Get()
  root(@Param('channelId') channelId: string) {
    return this.timersService.getList(channelId);
  }

  @UseGuards(JwtAuthGuard, DashboardAccessGuard)
  @Get(':id')
  async findOne(@Param('channelId') channelId: string, @Param('id') id: string) {
    const result = await this.timersService.findOne(channelId, id);
    return result;
  }

  @UseGuards(JwtAuthGuard, DashboardAccessGuard)
  @Post()
  async create(@Param('channelId') channelId: string, @Body() body: CreateTimerDto) {
    const result = await this.timersService.create(channelId, body);
    return result;
  }

  @UseGuards(JwtAuthGuard, DashboardAccessGuard)
  @Delete(':id')
  async delete(@Param('channelId') channelId: string, @Param('id') id: string) {
    const result = await this.timersService.delete(channelId, id);
    return result;
  }

  @UseGuards(JwtAuthGuard, DashboardAccessGuard)
  @Put(':id')
  async update(
    @Param('channelId') channelId: string,
    @Param('id') id: string,
    @Body() body: CreateTimerDto,
  ) {
    const result = await this.timersService.update(channelId, id, body);
    return result;
  }
}
