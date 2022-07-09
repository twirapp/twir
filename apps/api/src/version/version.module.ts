import { Module } from '@nestjs/common';

import { VersionController } from './version.controller.js';
import { VersionService } from './version.service.js';

@Module({
  controllers: [VersionController],
  providers: [VersionService],
})
export class VersionModule { }
