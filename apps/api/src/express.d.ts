import Twitch from '@nestjs-hybrid-auth/twitch';

import { JwtPayload } from '../http/jwt/jwt.strategy';


declare module 'express' {
  interface Request {
    hybridAuthResult: Twitch.TwitchAuthResult;
    user: JwtPayload & { iat: number; exp: number };
  }
}
