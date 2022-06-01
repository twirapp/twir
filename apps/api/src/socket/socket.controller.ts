import { CacheTTL, UseGuards, UseInterceptors } from '@nestjs/common';
import {
  MessageBody,
  SubscribeMessage,
  WebSocketGateway,
  WebSocketServer,
} from '@nestjs/websockets';
import { ClientToServerEvents, EventParams } from '@tsuwari/shared';
import { Server } from 'socket.io';

import { DashboardAccessGuard } from '../guards/DashboardAccess.guard.js';
import { WsJwtGuard } from '../guards/jwtws.guard.js';
import { CustomCacheInterceptor } from '../helpers/customCacheInterceptor.js';
import { SocketService } from './socket.service.js';

@WebSocketGateway({
  cors: {
    origin: '*',
  },
})
export class SocketController {
  @WebSocketServer()
  server: Server;

  constructor(private readonly socketService: SocketService) { }


  @CacheTTL(10)
  @UseInterceptors(CustomCacheInterceptor(ctx => {
    const req = ctx.switchToWs().getData() as EventParams<ClientToServerEvents, 'isBotMod'>[0];

    return `nest:cache:socket/isBotMod:${req.channelId}`;
  }))
  @UseGuards(WsJwtGuard, DashboardAccessGuard)
  @SubscribeMessage<keyof ClientToServerEvents>('isBotMod')
  isBotMod(@MessageBody() body: EventParams<ClientToServerEvents, 'isBotMod'>[0]) {
    return this.socketService.isBotMod(body);
  }

  @UseGuards(WsJwtGuard, DashboardAccessGuard)
  @SubscribeMessage<keyof ClientToServerEvents>('botJoin')
  botJoin(@MessageBody() body: EventParams<ClientToServerEvents, 'botJoin'>[0]) {
    return this.socketService.botJoinOrLeave('join', body);
  }

  @UseGuards(WsJwtGuard, DashboardAccessGuard)
  @SubscribeMessage<keyof ClientToServerEvents>('botPart')
  botPart(@MessageBody() body: EventParams<ClientToServerEvents, 'botPart'>[0]) {
    return this.socketService.botJoinOrLeave('part', body);
  }
}