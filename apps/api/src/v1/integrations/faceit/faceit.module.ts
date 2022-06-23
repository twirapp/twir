import { Module } from '@nestjs/common';

import { FaceitController } from './faceit.controller.js';
import { FaceitService } from './faceit.service.js';

@Module({
  controllers: [FaceitController],
  providers: [FaceitService],
})
export class FaceitModule { }
