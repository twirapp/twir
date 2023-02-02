import { Module } from '@nestjs/common';

import { EmotesService } from './emotes.service.js';

@Module({
  imports: [],
  providers: [EmotesService],
})
export class EmotesModule { }
