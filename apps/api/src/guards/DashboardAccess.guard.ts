import { CanActivate, ExecutionContext, HttpException, Injectable } from '@nestjs/common';
import { DashboardAccess } from '@tsuwari/typeorm/entities/DashboardAccess';
import { User } from '@tsuwari/typeorm/entities/User';
import Express from 'express';

import { typeorm } from '../index.js';

@Injectable()
export class DashboardAccessGuard implements CanActivate {
  async canActivate(context: ExecutionContext) {
    const request = context.switchToHttp().getRequest() as Express.Request;
    if (!request.params?.channelId || !request.user?.id)
      throw new HttpException('DashboardId not specified.', 400);

    if (request.user.id === request.params.channelId) {
      return true;
    }

    const [requestUser, dashBoardAccess] = await Promise.all([
      typeorm.getRepository(User).findOneBy({ id: request.user.id }),
      typeorm.getRepository(DashboardAccess).count({
        where: {
          channelId: request.params.channelId,
          userId: request.user.id,
        },
      }),
    ]);

    return requestUser?.isBotAdmin || !!dashBoardAccess;
  }
}
