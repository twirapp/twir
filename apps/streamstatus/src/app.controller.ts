import { Controller, Logger } from '@nestjs/common';
import { GrpcMethod, GrpcStreamCall, GrpcStreamMethod } from '@nestjs/microservices';
import { StreamStatus } from '@tsuwari/grpc';
import { Observable, Subject } from 'rxjs';

import { AppService } from './app.service.js';

@Controller()
export class AppController implements StreamStatus.Main {
  private readonly logger = new Logger('StreamStatus');

  constructor(private readonly appService: AppService) { }

  @GrpcStreamMethod('Main', 'CacheStreams')
  cacheStreams(data: Observable<StreamStatus.CacheStreamRequest>): Observable<StreamStatus.CachedStreamResult> {
    const subject = new Subject<StreamStatus.CachedStreamResult>();

    data.subscribe({
      next: (m) => {
        this.logger.log(`Starting to process ${m.channelsIds?.length} streams`);
        this.appService.handleChannels(m.channelsIds!).then(() => {
          subject.next({ success: true });
        }).catch(() => {
          subject.next({ success: false });
        });
      },
      complete: () => subject.complete(),
    });

    return subject.asObservable();
  }
}
