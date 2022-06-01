import { Module } from '@nestjs/common';

import { GreetingsController } from './greetings.controller.js';
import { GreetingsService } from './greetings.service.js';

@Module({
  controllers: [GreetingsController],
  providers: [GreetingsService],
})
export class GreetingsModule { }
