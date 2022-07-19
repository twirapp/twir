import { Body, Controller, Post, Req, UseGuards } from '@nestjs/common';
import Express from 'express';

import { JwtAuthGuard } from '../../jwt/jwt.guard.js';
import type { FeedBackPostDto } from './dto/post.dto.js';
import { FeedbackService } from './feedback.service.js';

@Controller('v1/feedback')
export class FeedbackController {
  constructor(private readonly service: FeedbackService) {

  }

  @Post()
  @UseGuards(JwtAuthGuard)
  postFeedBack(@Body() body: FeedBackPostDto, @Req() req: Express.Request) {
    return this.service.postFeedBack(body, req.user.id);
  }
}
