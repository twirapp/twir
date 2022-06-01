import { Module } from '@nestjs/common';

import { TimersController } from './timers.controller.js';

@Module({
  controllers: [TimersController],
})
export class TimersModule {}
