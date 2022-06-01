import { Module } from '@nestjs/common';

import { TimersController } from './timers.controller.js';
import { TimersService } from './timers.service.js';

@Module({
  controllers: [TimersController],
  providers: [TimersService],
})
export class TimersModule { }
