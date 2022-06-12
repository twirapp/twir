import { TwitchAuthModule } from '@nestjs-hybrid-auth/twitch';
import { Module } from '@nestjs/common';
import { config } from '@tsuwari/config';

import { JwtAuthModule } from '../jwt/jwt.module.js';
import { AuthController } from './auth.controller.js';
import { AuthService } from './auth.service.js';

export const scope = ['moderation:read', 'channel:manage:broadcast'];

@Module({
  imports: [
    TwitchAuthModule.forRoot({
      clientID: config.TWITCH_CLIENTID,
      clientSecret: config.TWITCH_CLIENTSECRET,
      callbackURL: config.TWITCH_CALLBACKURL,
      forceVerify: false,
      scope,
    } as any),
    JwtAuthModule,
  ],
  controllers: [AuthController],
  providers: [AuthService],
})
export class AuthModule { }
