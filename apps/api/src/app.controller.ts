import { Controller, Get, Res } from '@nestjs/common';
import { Response } from 'express';

import { AppService } from './app.service.js';
import { prometheus } from './prometheus.js';

@Controller()
export class AppController {
  constructor(private readonly appService: AppService) { }

  /* @Get('/')
  root() {
    return '';
  } */

  @Get('/metrics')
  async metrics(@Res() res: Response) {
    res.contentType(prometheus.contentType);
    res.send(await prometheus.register.metrics());
  }
}

