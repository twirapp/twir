/* eslint-disable import/no-cycle */
import {
  Column,
  Entity,
  OneToMany,
  PrimaryColumn,
} from 'typeorm';

import { ChannelIntegration } from './ChannelIntegration';

export enum IntegrationService {
  LASTFM = 'LASTFM',
  VK = 'VK',
  FACEIT = 'FACEIT',
  SPOTIFY = 'SPOTIFY',
  DONATIONALERTS = 'DONATIONALERTS',
  STREAMLABS = 'STREAMLABS',
}

@Entity('integrations', { schema: 'public' })
export class Integration {
  @PrimaryColumn('text', {
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

  @OneToMany(() => ChannelIntegration, _ => _.integration)
  channelsIntegrations?: ChannelIntegration[];
}
