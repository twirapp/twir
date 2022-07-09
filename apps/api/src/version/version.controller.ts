import { CacheKey, CacheTTL, Controller, Get } from '@nestjs/common';

import { VersionService } from './version.service.js';

@Controller('version')
export class VersionController {
  constructor(private readonly service: VersionService) {

  }

  @CacheKey('api/version')
  @CacheTTL(20)
  @Get()
  version() {
    return this.service.getCommitSha();
  }
}
