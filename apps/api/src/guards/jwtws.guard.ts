import { CanActivate, ExecutionContext, Injectable } from '@nestjs/common';
import { config } from '@tsuwari/config';
import jwt from 'jsonwebtoken';
import { Socket } from 'socket.io';


@Injectable()
export class WsJwtGuard implements CanActivate {
  async canActivate(context: ExecutionContext) {
    const client = context.switchToWs().getClient() as Socket;
    const authToken = client.handshake.auth.token;

    if (!authToken) {
      return false;
    }

    try {
      const verify = jwt.verify(authToken, config.JWT_ACCESS_SECRET);
      return !!verify;
    } catch (e) {
      if (e instanceof jwt.TokenExpiredError) {
        if (e.message.includes('jwt expired')) return false;
      }

      return false;
    }

  }
}