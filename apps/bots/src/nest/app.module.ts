import { Module } from '@nestjs/common';

import { AppController } from './app.controller.js';
import { ParserModule } from './parser/parser.module.js';
import { TimersModule } from './timers/timers.module.js';

@Module({
  imports: [TimersModule, ParserModule],
  controllers: [AppController],
  providers: [],
})
export class AppModule { }
