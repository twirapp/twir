import { Module } from '@nestjs/common';

import { DashboardAccessModule } from './dashboard-access/dashboard-access.module.js';

@Module({
  imports: [DashboardAccessModule],
})
export class SettingsModule { }
