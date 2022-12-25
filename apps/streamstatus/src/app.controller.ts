import { Controller, Logger, OnModuleInit } from '@nestjs/common';

import { AppService } from './app.service.js';
import { pubSub } from './pubsub.js';

@Controller()
export class AppController implements OnModuleInit {
  constructor(private readonly appService: AppService) {}

  async onModuleInit() {
    pubSub.subscribe('streams.online', (data) => {
      const parsedData = JSON.parse(data);
      this.appService.handleOnline(parsedData);
    });
    pubSub.subscribe('streams.offline', (data) => {
      const parsedData = JSON.parse(data);
      this.appService.handleOffline(parsedData);
    });
    pubSub.subscribe('stream.update', (data) => {
      const parsedData = JSON.parse(data);
      this.appService.handleUpdate(parsedData);
    });
  }
}
