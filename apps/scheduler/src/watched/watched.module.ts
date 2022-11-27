import { Module } from '@nestjs/common';

import { WatchedService } from './watched.service.js';

@Module({
  imports: [],
  controllers: [],
  providers: [WatchedService],
})
export class WatchedModule {}
