import { Module } from '@nestjs/common';

import { StreamStatusService } from './statuschecker.service.js';

@Module({
  imports: [],
  providers: [StreamStatusService],
})
export class StreamStatusModule { }
