import { Injectable, CanActivate, ExecutionContext } from '@nestjs/common';
import { PrismaService } from '@tsuwari/prisma';
import Express from 'express';

@Injectable()
export class IsAdminGuard implements CanActivate {
  constructor(private readonly prisma: PrismaService) { }

  async canActivate(
    context: ExecutionContext,
  ) {
    const request = context.switchToHttp().getRequest() as Express.Request;
    if (!request.user?.id) return false;
    const user = await this.prisma.user.findFirst({ where: { id: request.user.id, isBotAdmin: true } });

    if (!user) return false;
    return true;
  }
}