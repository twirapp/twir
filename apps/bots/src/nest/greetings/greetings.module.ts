import { Module } from '@nestjs/common';

import { GreetingsController } from './greetings.controller.js';

@Module({
  controllers: [GreetingsController],
})
export class GreetingsModule {}
