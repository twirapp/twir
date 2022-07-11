import { Module } from '@nestjs/common';

import { IncreaseWatchedService } from './increase-watched.service.js';

@Module({
  providers: [IncreaseWatchedService],
})
export class IncreaseWatchedModule { }
