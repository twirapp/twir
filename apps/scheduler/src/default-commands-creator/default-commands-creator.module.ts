import { Module } from '@nestjs/common';

import { DefaultCommandsCreatorService } from './default-commands-creator.service.js';

@Module({
  controllers: [],
  imports: [],
  providers: [DefaultCommandsCreatorService],
})
export class DefaultCommandsCreatorModule {}
