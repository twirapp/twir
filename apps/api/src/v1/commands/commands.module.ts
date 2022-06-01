import { Module } from '@nestjs/common';

import { CommandsController } from './commands.controller.js';
import { CommandsService } from './commands.service.js';

@Module({
  controllers: [CommandsController],
  providers: [CommandsService],
})
export class CommandsModule { }
