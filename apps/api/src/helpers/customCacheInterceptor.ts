import { Injectable, CacheInterceptor, ExecutionContext, mixin } from '@nestjs/common';

export function CustomCacheInterceptor(fn: (ctx: ExecutionContext) => string | undefined) {
  @Injectable()
  class Interceptor extends CacheInterceptor {
    trackBy(ctx: ExecutionContext) {
      const result = fn(ctx);
      return result;
    }
  }

  return mixin(Interceptor as any);
}
