import { Module } from '@nestjs/common';

import { VariablesController } from './variables.controller.js';
import { VariablesService } from './variables.service.js';

@Module({
  controllers: [VariablesController],
  providers: [VariablesService],
})
export class VariablesModule { }
