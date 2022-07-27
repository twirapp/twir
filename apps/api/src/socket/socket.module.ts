import { Module } from '@nestjs/common';

import { SocketController } from './socket.controller.js';
import { SocketService } from './socket.service.js';


@Module({
  controllers: [],
  providers: [SocketController, SocketService],
})
export class SocketModule { }
