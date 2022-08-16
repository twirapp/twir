import { CanActivate, ExecutionContext, HttpException, Injectable } from '@nestjs/common';
import { PrismaService } from '@tsuwari/prisma';
import Express from 'express';

@Injectable()
export class DashboardAccessGuard implements CanActivate {
  constructor(private readonly prisma: PrismaService) {}

  async canActivate(context: ExecutionContext) {
    const request = context.switchToHttp().getRequest() as Express.Request;
    if (!request.params?.channelId || !request.user?.id)
      throw new HttpException('DashboardId not specified.', 400);

    if (request.user.id === request.params.channelId) {
      return true;
    }

    const [requestUser, dashBoardAccess] = await Promise.all([
      this.prisma.user.findFirst({ where: { id: request.user.id } }),
      this.prisma.dashboardAccess.count({
        where: {
          channelId: request.params.channelId,
          userId: request.user.id,
        },
      }),
    ]);

    return requestUser?.isBotAdmin || !!dashBoardAccess;
  }
}
