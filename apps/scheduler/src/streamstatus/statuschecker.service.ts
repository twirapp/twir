import { Inject, Injectable } from '@nestjs/common';
import { Interval } from '@nestjs/schedule';
import { config } from '@tsuwari/config';
import { ClientProxy } from '@tsuwari/shared';
import { Channel } from '@tsuwari/typeorm/entities/Channel';
import _ from 'lodash';

import { typeorm } from '../index.js';

@Injectable()
export class StreamStatusService {
  constructor(@Inject('NATS') private nats: ClientProxy) {}

  @Interval('streamstatus', config.isDev ? 10000 : 5 * 60 * 1000)
  async updateStatuses() {
    const channelsIds = await typeorm.getRepository(Channel).find({
      where: { isEnabled: true },
      select: { id: true },
    });

    const chunks = _.chunk(
      channelsIds.map((c) => c.id),
      100,
    );

    for (const channelsIds of chunks) {
      await this.nats.send('streamstatuses.process', channelsIds).toPromise();
    }
  }
}
