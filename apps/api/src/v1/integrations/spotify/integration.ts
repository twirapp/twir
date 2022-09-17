import { Injectable } from '@nestjs/common';
import { SpotifyIntegration } from '@tsuwari/spotify';
import { Integration } from '@tsuwari/typeorm/entities/Integration';

import { typeorm } from '../../../index.js';

@Injectable()
export class SpotifyIntegrationService extends SpotifyIntegration {
  constructor() {
    super(typeorm.getRepository(Integration));
  }
}
