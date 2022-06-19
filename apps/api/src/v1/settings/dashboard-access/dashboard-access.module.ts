import { Module } from '@nestjs/common';

import { DashboardAccessController } from './dashboard-access.controller.js';
import { DashboardAccessService } from './dashboard-access.service.js';

@Module({
  providers: [DashboardAccessService],
  controllers: [DashboardAccessController],
})
export class DashboardAccessModule { }
