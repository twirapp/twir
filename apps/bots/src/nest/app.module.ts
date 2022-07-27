import { Module } from '@nestjs/common';
import { RedisORMModule } from '@tsuwari/redis';

import { AppController } from './app.controller.js';
import { CacherModule } from './cacher/cacher.module.js';
import { ParserModule } from './parser/parser.module.js';
import { TimersModule } from './timers/timers.module.js';

@Module({
  imports: [RedisORMModule, TimersModule, ParserModule, CacherModule],
  controllers: [AppController],
  providers: [],
})
export class AppModule { }
