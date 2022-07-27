import { Global, Injectable } from '@nestjs/common';
import { Reflector } from '@nestjs/core';
import { InjectThrottlerOptions, InjectThrottlerStorage, ThrottlerGuard as OriginalGuard, ThrottlerStorage } from '@nestjs/throttler';
import type { ThrottlerModuleOptions } from '@nestjs/throttler';
import Express from 'express';

@Global()
@Injectable()
export class ThrottlerGuard extends OriginalGuard {
  headerPrefix = '';
  errorMessage = 'Too many requests.';

  constructor(
    @InjectThrottlerOptions()
    protected readonly options: ThrottlerModuleOptions,
    @InjectThrottlerStorage()
    storage: ThrottlerStorage,
    protected readonly reflector: Reflector) {
    super(options, storage, reflector);
  }

  protected getTracker(req: Express.Request): string {
    return req.user?.id;
  }
}
