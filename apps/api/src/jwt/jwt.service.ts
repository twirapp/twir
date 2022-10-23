import { Injectable } from '@nestjs/common';
import { JwtService } from '@nestjs/jwt';
import { config } from '@tsuwari/config';

import { JwtPayload } from './jwt.strategy.js';

export interface Payload {
  username: string;
  id: string;
  scopes: string[];
}

@Injectable()
export class JwtAuthService {
  constructor(private jwtService: JwtService) {}

  private async signAccessToken(payload: Payload) {
    return await this.jwtService.signAsync(payload, {
      expiresIn: config.JWT_EXPIRES_IN,
      secret: config.JWT_ACCESS_SECRET,
    });
  }

  private async signRefreshToken(payload: Payload) {
    return await this.jwtService.signAsync(payload, {
      expiresIn: '31d',
      secret: config.JWT_REFRESH_SECRET,
    });
  }

  private async verifyRefreshToken(refreshToken: string) {
    return await this.jwtService.verifyAsync<JwtPayload>(refreshToken, {
      secret: config.JWT_REFRESH_SECRET,
    });
  }

  async refreshAccessToken(refreshToken: string) {
    const payload = await this.verifyRefreshToken(refreshToken);
    return await this.signAccessToken({
      id: payload.id,
      scopes: payload.scopes,
      username: payload.login,
    });
  }

  async generateKeypair(payload: Payload) {
    const accessToken = await this.signAccessToken(payload);
    const refreshToken = await this.signRefreshToken(payload);

    return { accessToken, refreshToken };
  }
}
