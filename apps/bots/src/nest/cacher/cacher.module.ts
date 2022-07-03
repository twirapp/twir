import { Module } from '@nestjs/common';

import { CacherService } from './cacher.service.js';

@Module({
  providers: [CacherService],
})
export class CacherModule { }
