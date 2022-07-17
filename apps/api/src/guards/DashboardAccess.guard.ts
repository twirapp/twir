import { Injectable, CanActivate, ExecutionContext, HttpException } from '@nestjs/common';
import { PrismaService } from '@tsuwari/prisma';
import Express from 'express';
import { Socket } from 'socket.io';


@Injectable()
export class DashboardAccessGuard implements CanActivate {
  constructor(private readonly prisma: PrismaService) { }

  async canActivate(
    context: ExecutionContext,
  ) {
    const type = context.getType();

    if (type === 'ws') {
      return this.#handleWs(context);
    }

    const request = context.switchToHttp().getRequest() as Express.Request;
    if (!request.params?.channelId || !request.user?.id) throw new HttpException('DashboardId not specified.', 400);

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

  async #handleWs(context: ExecutionContext) {
    const client = context.switchToWs().getClient() as Socket;
    const query = client.handshake.query;
    if (!query?.channelId || !query?.userId) throw new Error('DashboardId not specified.');

    if (Array.isArray(query.userId) || Array.isArray(query.channelId)) throw new Error('UserId or channelId cannot be array.');
    if (query.channelId === query.userId) {
      return true;
    }

    const [requestUser, dashBoardAccess] = await Promise.all([
      this.prisma.user.findFirst({ where: { id: query.userId } }),
      this.prisma.dashboardAccess.count({
        where: {
          channelId: query.channelId,
          userId: query.userId,
        },
      }),
    ]);

    return requestUser?.isBotAdmin || !!dashBoardAccess;
  }
}