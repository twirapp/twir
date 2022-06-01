import { TwitchAuthResult } from '@nestjs-hybrid-auth/twitch';

import 'express';
import { JwtPayload } from '../http/jwt/jwt.strategy';

declare module 'express' {
  interface Request {
    hybridAuthResult: TwitchAuthResult;
    user: JwtPayload & { iat: number; exp: number };
  }
}
