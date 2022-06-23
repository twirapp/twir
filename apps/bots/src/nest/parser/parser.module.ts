import { Module } from '@nestjs/common';

import { ParserService } from './parser.service.js';

@Module({
  providers: [ParserService],
})
export class ParserModule { }
