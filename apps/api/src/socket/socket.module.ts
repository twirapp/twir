import { Module } from '@nestjs/common';

import { RedisService } from '../redis.service.js';
import { SocketController } from './socket.controller.js';
import { SocketService } from './socket.service.js';


@Module({
  controllers: [],
  providers: [RedisService, SocketController, SocketService],
})
export class SocketModule { }
