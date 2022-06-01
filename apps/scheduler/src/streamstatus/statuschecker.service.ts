import { Inject, Injectable, OnModuleInit } from '@nestjs/common';
import { ClientGrpc } from '@nestjs/microservices';
import { Interval } from '@nestjs/schedule';
import { config } from '@tsuwari/config';
import { StreamStatus } from '@tsuwari/grpc';
import { PrismaService } from '@tsuwari/prisma';
import _ from 'lodash';
import { ReplaySubject } from 'rxjs';

@Injectable()
export class StreamStatusService implements OnModuleInit {
  private statusMicroservice: StreamStatus.Main;

  constructor(
    @Inject('STREAMSTATUS_MICROSERVICE') private client: ClientGrpc,
    private readonly prisma: PrismaService,
  ) { }

  onModuleInit(): void {
    this.statusMicroservice = this.client.getService<StreamStatus.Main>('Main');
  }

  @Interval(config.isDev ? 10000 : 5 * 60 * 1000)
  async updateStatuses() {
    const channelsIds = await this.prisma.channel.findMany({
      select: {
        id: true,
      },
    });
    const helloRequest$ = new ReplaySubject<StreamStatus.CacheStreamRequest>();

    const chunks = _.chunk(channelsIds.map(c => c.id), 100);

    for (const channelsIds of chunks) {
      helloRequest$.next({ channelsIds });
    }
    helloRequest$.complete();

    this.statusMicroservice.cacheStreams(helloRequest$.asObservable()).subscribe();
  }
}

