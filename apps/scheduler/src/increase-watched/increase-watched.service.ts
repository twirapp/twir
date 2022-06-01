import { Inject, Injectable } from '@nestjs/common';
import { ClientGrpc } from '@nestjs/microservices';
import { Interval } from '@nestjs/schedule';
import { config } from '@tsuwari/config';
import { Watched } from '@tsuwari/grpc';
import _ from 'lodash';
import { ReplaySubject } from 'rxjs';

import { RedisService } from '../redis.service.js';

@Injectable()
export class IncreaseWatchedService {
  private watchedMicroservice: Watched.Main;

  constructor(
    @Inject('WATCHED_MICROSERVICE') private client: ClientGrpc,
    private readonly redis: RedisService,
  ) { }

  onModuleInit(): void {
    this.watchedMicroservice = this.client.getService<Watched.Main>('Main');
  }

  @Interval(config.isDev ? 10000 : 5 * 60 * 1000)
  async updateStatuses() {
    const streamsKeys = await this.redis.keys('streams:*');
    if (!streamsKeys.length) return;

    const helloRequest$ = new ReplaySubject<Watched.IncreaseWatchedRequest>();
    const chunks = _.chunk(streamsKeys, 100);

    for (const channelsIds of chunks) {
      helloRequest$.next({ streamsIds: channelsIds });
    }
    helloRequest$.complete();

    this.watchedMicroservice.increaseWatched(helloRequest$.asObservable()).subscribe();
  }
}
