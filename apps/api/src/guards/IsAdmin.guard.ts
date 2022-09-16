import { CanActivate, ExecutionContext, Injectable } from '@nestjs/common';
import { User } from '@tsuwari/typeorm/entities/User';
import Express from 'express';

import { typeorm } from '../index.js';

@Injectable()
export class IsAdminGuard implements CanActivate {
  async canActivate(context: ExecutionContext) {
    const request = context.switchToHttp().getRequest() as Express.Request;
    if (!request.user?.id) return false;
    const user = await typeorm
      .getRepository(User)
      .findOneBy({ id: request.user.id, isBotAdmin: true });

    if (!user) return false;
    return true;
  }
}
