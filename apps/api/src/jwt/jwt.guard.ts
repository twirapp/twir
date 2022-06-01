import { Injectable, UnauthorizedException } from '@nestjs/common';
import { AuthGuard } from '@nestjs/passport';
import jwt from 'jsonwebtoken';

@Injectable()
export class JwtAuthGuard extends AuthGuard('jwt') {
  handleRequest(err: any, user: any, info: any, context: any, status: any) {
    if (info instanceof jwt.JsonWebTokenError) {
      throw new UnauthorizedException('Invalid Token!');
    }

    return super.handleRequest(err, user, info, context, status);
  }
}

