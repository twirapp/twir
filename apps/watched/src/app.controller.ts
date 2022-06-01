import { Controller, Logger } from '@nestjs/common';
import { GrpcStreamMethod } from '@nestjs/microservices';
import { Watched } from '@tsuwari/grpc';
import { Observable, Subject } from 'rxjs';

import { AppService } from './app.service.js';

@Controller()
export class AppController implements Watched.Main {
  private readonly logger = new Logger('Watcher');

  constructor(private readonly appService: AppService) { }

  @GrpcStreamMethod('Main', 'IncreaseWatched')
  increaseWatched(data: Observable<Watched.IncreaseWatchedRequest>): Observable<Watched.IncreaseWatchedResult> {
    const subject = new Subject<Watched.IncreaseWatchedResult>();

    data.subscribe({
      next: (m) => {
        this.logger.log(`Starting to process ${m.streamsIds?.length} streams`);
        this.appService.handleIncrease(m.streamsIds!).then(() => {
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
