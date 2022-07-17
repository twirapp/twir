import { Controller, Get, Res } from '@nestjs/common';
import Express from 'express';

import { AppService } from './app.service.js';
import { prometheus } from './prometheus.js';

@Controller()
export class AppController {
  constructor(private readonly appService: AppService) { }

  @Get('/metrics')
  async metrics(@Res() res: Express.Response) {
    res.contentType(prometheus.contentType);
    res.send(await prometheus.register.metrics());
  }
}

