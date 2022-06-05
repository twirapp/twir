import { Column, Entity, OneToMany, PrimaryGeneratedColumn } from 'typeorm';

import { ChannelIntegration } from './ChannelIntegration.js';

export enum IntegrationService {
  LASTFM = 'lastfm',
  VK = 'vk',
  FACEIT = 'faceit',
  SPOTIFY = 'spotify'
}

@Entity('integrations')
export class Integration {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @Column({
    type: 'enum',
    enum: IntegrationService,
  })
  service: IntegrationService;

  @Column()
  accessToken?: string;

  @Column()
  refreshToken?: string;

  @Column()
  clientId?: string;

  @Column()
  clientSecret?: string;

  @Column()
  redirectUrl?: string;

  @OneToMany(() => ChannelIntegration, i => i.integration)
  channelIntegrations: ChannelIntegration[];
}