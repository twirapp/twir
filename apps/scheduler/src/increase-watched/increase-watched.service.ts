import { Injectable } from '@nestjs/common';
import { RedisService } from '@tsuwari/shared';
import _ from 'lodash';

@Injectable()
export class IncreaseWatchedService {

  constructor(
    private readonly redis: RedisService,
  ) { }

  // @Timeout(500)
  /* async updateStatuses() {
    console.log('start');
    const streamsKeys = await this.redis.keys('streams:*');
    if (!streamsKeys.length) return;

    const helloRequest$ = new ReplaySubject<Watched.IncreaseWatchedRequest>();
    const chunks = _.chunk(streamsKeys, 100);

    for (const channelsIds of chunks) {
      helloRequest$.next({ streamsIds: channelsIds });
    }
    helloRequest$.complete();

    this.watchedMicroservice.increaseWatched(helloRequest$.asObservable()).subscribe();
  } */
}
