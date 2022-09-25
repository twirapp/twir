import { IsBoolean, IsOptional } from 'class-validator';

export class UpdateStreamlabsIntegrationDto {
  @IsBoolean()
  @IsOptional()
  enabled?: boolean;
}
