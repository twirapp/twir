import { Injectable } from '@nestjs/common';
import { JwtService } from '@nestjs/jwt';
import { config } from '@tsuwari/config';

import { JwtPayload } from './jwt.strategy.js';

@Injectable()
export class JwtAuthService {
  constructor(private jwtService: JwtService) { }

  login(user: JwtPayload) {
    const payload = { username: user.login, id: user.id };
    return {
      accessToken: this.jwtService.sign(payload, {
        expiresIn: config.JWT_EXPIRES_IN,
        secret: config.JWT_ACCESS_SECRET,
      }),
      refreshToken: this.jwtService.sign(payload, { expiresIn: '31d', secret: config.JWT_REFRESH_SECRET }),
    };
  }

  async refresh(token: string) {
    const user = await this.jwtService.verifyAsync<JwtPayload>(token, { secret: config.JWT_REFRESH_SECRET });
    return this.login(user);
  }
}
