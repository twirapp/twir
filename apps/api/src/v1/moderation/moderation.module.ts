import { Module } from '@nestjs/common';

import { ModerationController } from './moderation.controller.js';
import { ModerationService } from './moderation.service.js';

@Module({
  controllers: [ModerationController],
  providers: [ModerationService],
})
export class ModerationModule { }
