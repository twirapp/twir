import { Module } from '@nestjs/common';

import { AppController } from './app.controller.js';
import { TimersModule } from './timers/timers.module.js';

@Module({
  imports: [TimersModule],
  controllers: [AppController],
  providers: [],
})
export class AppModule { }
