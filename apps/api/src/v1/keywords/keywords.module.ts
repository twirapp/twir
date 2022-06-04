import { Module } from '@nestjs/common';

import { KeywordsController } from './keywords.controller.js';
import { KeywordsService } from './keywords.service.js';

@Module({
  providers: [KeywordsService],
  controllers: [KeywordsController],
})
export class KeywordsModule { }
