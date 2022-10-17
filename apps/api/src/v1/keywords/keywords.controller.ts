import { Body, Controller, Delete, Get, Param, Patch, Post, UseGuards } from '@nestjs/common';

import { DashboardAccessGuard } from '../../guards/DashboardAccess.guard.js';
import { JwtAuthGuard } from '../../jwt/jwt.guard.js';
import { CreateKeywordDto } from './dto/create.js';
import { KeywordsService } from './keywords.service.js';

@Controller('v1/channels/:channelId/keywords')
export class KeywordsController {
  constructor(private readonly keywordsService: KeywordsService) {}

  @UseGuards(JwtAuthGuard, DashboardAccessGuard)
  @Get()
  root(@Param('channelId') channelId: string) {
    return this.keywordsService.getList(channelId);
  }

  @UseGuards(JwtAuthGuard, DashboardAccessGuard)
  @Delete(':keywordId')
  async delete(@Param('channelId') channelId: string, @Param('keywordId') keywordId: string) {
    const result = await this.keywordsService.delete(channelId, keywordId);
    return result;
  }

  @UseGuards(JwtAuthGuard, DashboardAccessGuard)
  @Post()
  async create(@Param('channelId') channelId: string, @Body() body: CreateKeywordDto) {
    const result = await this.keywordsService.create(channelId, body);
    return result;
  }

  @UseGuards(JwtAuthGuard, DashboardAccessGuard)
  @Patch(':keywordId')
  async patch(
    @Param('channelId') channelId: string,
    @Param('keywordId') keywordId: string,
    @Body() body: CreateKeywordDto,
  ) {
    const result = await this.keywordsService.patch(channelId, keywordId, body);
    return result;
  }
}
