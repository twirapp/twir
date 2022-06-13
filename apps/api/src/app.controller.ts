import { Controller, Get, Res } from '@nestjs/common';
import { EventPattern, Payload } from '@nestjs/microservices';
import { Response } from 'express';

import { AppService } from './app.service.js';
import { prometheus } from './prometheus.js';

@Controller()
export class AppController {
  constructor(private readonly appService: AppService) { }

  @Get('/metrics')
  async metrics(@Res() res: Response) {
    res.contentType(prometheus.contentType);
    res.send(await prometheus.register.metrics());
  }
}

