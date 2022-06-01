import { Module } from '@nestjs/common';

import { CommandsController } from './commands.controller.js';

@Module({
  controllers: [CommandsController],
})
export class CommandsModule {}
