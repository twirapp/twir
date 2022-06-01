import { Module } from '@nestjs/common';

import { AppController } from './app.controller.js';
import { CommandsModule } from './commands/commands.module.js';
import { GreetingsModule } from './greetings/greetings.module.js';
import { TimersModule } from './timers/timers.module.js';

@Module({
  imports: [CommandsModule, GreetingsModule, TimersModule],
  controllers: [AppController],
  providers: [],
})
export class AppModule {}
