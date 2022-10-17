import { Body, Controller, Delete, Get, Param, Post, Put, UseGuards } from '@nestjs/common';

import { DashboardAccessGuard } from '../../guards/DashboardAccess.guard.js';
import { JwtAuthGuard } from '../../jwt/jwt.guard.js';
import { CreateVariableDto } from './dto/create.js';
import { VariablesService } from './variables.service.js';

@Controller('v1/channels/:channelId/variables')
export class VariablesController {
  constructor(private readonly variablesService: VariablesService) {}

  @Get('builtin')
  @UseGuards(JwtAuthGuard, DashboardAccessGuard)
  builtIn() {
    return this.variablesService.getBuildInVariables();
  }

  @UseGuards(JwtAuthGuard, DashboardAccessGuard)
  @Get()
  root(@Param('channelId') channelId: string) {
    return this.variablesService.getList(channelId);
  }

  @UseGuards(JwtAuthGuard, DashboardAccessGuard)
  @Post()
  async create(@Param('channelId') channelId: string, @Body() body: CreateVariableDto) {
    const result = await this.variablesService.create(channelId, body);
    return result;
  }

  @UseGuards(JwtAuthGuard, DashboardAccessGuard)
  @Delete(':id')
  async delete(@Param('channelId') channelId: string, @Param('id') id: string) {
    const result = await this.variablesService.delete(channelId, id);
    return result;
  }

  @UseGuards(JwtAuthGuard, DashboardAccessGuard)
  @Put(':id')
  async update(
    @Param('channelId') channelId: string,
    @Param('id') id: string,
    @Body() body: CreateVariableDto,
  ) {
    const result = await this.variablesService.update(channelId, id, body);
    return result;
  }
}
