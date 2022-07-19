import { Injectable } from '@nestjs/common';
import { ThrottlerGuard as OriginalGuard } from '@nestjs/throttler';
import Express from 'express';

@Injectable()
export class ThrottlerGuard extends OriginalGuard {
  protected getTracker(req: Express.Request): string {
    return req.user.id;
  }
}
