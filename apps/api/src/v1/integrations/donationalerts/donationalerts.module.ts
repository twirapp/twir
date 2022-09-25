import { Module } from '@nestjs/common';

import { DonationAlertsController } from './donationalerts.controller.js';
import { DonationAlertsService } from './donationalerts.service.js';

@Module({
  controllers: [DonationAlertsController],
  providers: [DonationAlertsService],
})
export class DonationAlertsModule {}
