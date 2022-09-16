/* eslint-disable import/no-cycle */
import { Column, Entity, Index, OneToMany } from 'typeorm';

import { ChannelIntegration } from './ChannelIntegration.js';

export enum IntegrationService {
  LASTFM = 'LASTFM',
  VK = 'VK',
  FACEIT = 'FACEIT',
  SPOTIFY = 'SPOTIFY',
  DONATIONALERTS = 'DONATIONALERTS',
}

@Index('integrations_pkey', ['id'], { unique: true })
@Entity('integrations', { schema: 'public' })
export class Integration {
  @Column('text', {
    primary: true,
    name: 'id',
    default: () => 'gen_random_uuid()',
  })
  id: string;

  @Column('enum', {
    name: 'service',
    enum: IntegrationService,
  })
  service: IntegrationService;

  @Column('text', { name: 'accessToken', nullable: true })
  accessToken: string | null;

  @Column('text', { name: 'refreshToken', nullable: true })
  refreshToken: string | null;

  @Column('text', { name: 'clientId', nullable: true })
  clientId: string | null;

  @Column('text', { name: 'clientSecret', nullable: true })
  clientSecret: string | null;

  @Column('text', { name: 'apiKey', nullable: true })
  apiKey: string | null;

  @Column('text', { name: 'redirectUrl', nullable: true })
  redirectUrl: string | null;

  @OneToMany(() => ChannelIntegration, (channelIntegration) => channelIntegration.integration)
  channelsIntegrations: ChannelIntegration[];
}
