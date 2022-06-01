import { Global, Module } from '@nestjs/common';

import { PrismaService } from './prisma.service.js';

@Global()
@Module({
  controllers: [],
  providers: [PrismaService],
  exports: [PrismaService],
})
export class PrismaModule { }