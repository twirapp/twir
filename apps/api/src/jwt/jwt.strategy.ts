import { Injectable } from '@nestjs/common';
import { PassportStrategy } from '@nestjs/passport';
import { config } from '@tsuwari/config';
import { Request } from 'express';
import { ExtractJwt, Strategy } from 'passport-jwt';


export type JwtPayload = {
  id: string;
  login: string;
  scopes: string[]
};

@Injectable()
export class JwtAuthStrategy extends PassportStrategy(Strategy) {
  constructor() {
    super({
      jwtFromRequest: (req: Request) => {
        const token = ExtractJwt.fromAuthHeaderAsBearerToken()(req);
        return token;
      },
      ignoreExpiration: false,
      secretOrKey: config.JWT_ACCESS_SECRET,
    });
  }

  async validate(payload: JwtPayload) {
    //TODO: check blacklist
    return payload;
  }
}
