import { IsBoolean, IsOptional } from 'class-validator';

export class UpdateSpotifyIntegrationDto {
  @IsBoolean()
  @IsOptional()
  enabled?: boolean;
}