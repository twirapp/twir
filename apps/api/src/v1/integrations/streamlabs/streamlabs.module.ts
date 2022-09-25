import { Module } from '@nestjs/common';

import { StreamlabsController } from './streamlabs.controller.js';
import { StreamlabsService } from './streamlabs.service.js';

@Module({
  controllers: [StreamlabsController],
  providers: [StreamlabsService],
})
export class StreamlabsModule {}
