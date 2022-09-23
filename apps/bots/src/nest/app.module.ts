import { Module } from '@nestjs/common';
import { RedisORMModule } from '@tsuwari/redis';

import { AppController } from './app.controller.js';
import { CacherModule } from './cacher/cacher.module.js';
import { ParserModule } from './parser/parser.module.js';

@Module({
  imports: [RedisORMModule, ParserModule, CacherModule],
  controllers: [AppController],
  providers: [],
})
export class AppModule {}
