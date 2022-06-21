import { Module } from '@nestjs/common';

import { VkController } from './vk.controller.js';
import { VkService } from './vk.service.js';

@Module({
  controllers: [VkController],
  providers: [VkService],
})
export class VkModule { }
