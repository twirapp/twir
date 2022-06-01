import { Module } from '@nestjs/common';
import { JwtModule } from '@nestjs/jwt';
import { config } from '@tsuwari/config';

import { JwtAuthService } from './jwt.service.js';
import { JwtAuthStrategy } from './jwt.strategy.js';

@Module({
  imports: [
    JwtModule.register({
      secret: config.JWT_ACCESS_SECRET,
      signOptions: {
        expiresIn: config.JWT_EXPIRES_IN,
      },
    }),
  ],
  providers: [JwtAuthStrategy, JwtAuthService],
  exports: [JwtModule, JwtAuthService],
})
export class JwtAuthModule { }
