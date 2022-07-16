import { Inject, Injectable } from '@nestjs/common';
import { Interval } from '@nestjs/schedule';
import { config } from '@tsuwari/config';
import { PrismaService } from '@tsuwari/prisma';
import { ClientProxy } from '@tsuwari/shared';
import _ from 'lodash';


@Injectable()
export class StreamStatusService {
  /* @Client({ transport: Transport.NATS, options: { servers: [config.NATS_URL] } })
  nats: ClientProxy; */

  constructor(
    private readonly prisma: PrismaService,
    @Inject('NATS') private nats: ClientProxy,
  ) { }

  @Interval('streamstatus', config.isDev ? 10000 : 5 * 60 * 1000)
  async updateStatuses() {
    const channelsIds = await this.prisma.channel.findMany({
      where: {
        isEnabled: true,
      },
      select: {
        id: true,
      },
    });

    const chunks = _.chunk(channelsIds.map(c => c.id), 100);

    for (const channelsIds of chunks) {
      await this.nats.send('streamstatuses.process', channelsIds).toPromise();
    }
  }
}

