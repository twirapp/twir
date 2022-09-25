import { IsBoolean, IsOptional } from 'class-validator';

export class UpdateDonationAlertsIntegrationDto {
  @IsBoolean()
  @IsOptional()
  enabled?: boolean;
}
