import { Inject, Injectable, Logger } from '@nestjs/common';
import { Interval } from '@nestjs/schedule';
import { config } from '@tsuwari/config';
import { ClientProxy } from '@tsuwari/shared';
import { ChannelDotaAccount } from '@tsuwari/typeorm/entities/ChannelDotaAccount';
import _ from 'lodash';

import { typeorm } from '../index.js';

@Injectable()
export class DotaService {
  #logger = new Logger(DotaService.name);

  constructor(@Inject('NATS') private nats: ClientProxy) {}

  @Interval('dota', config.isDev ? 10000 : 1 * 60 * 1000)
  async cacheDota() {
    const accounts = await typeorm.getRepository(ChannelDotaAccount).findBy({
      channel: {
        isEnabled: true,
      },
    });

    const chunks = _.chunk(
      accounts.map((a) => a.id),
      50,
    );

    for (const chunk of chunks) {
      await this.nats.emit('dota.cacheAccountsMatches', chunk).toPromise();
    }
  }
}
