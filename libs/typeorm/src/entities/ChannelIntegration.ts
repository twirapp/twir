/* eslint-disable import/no-cycle */
import { Column, Entity, Index, JoinColumn, ManyToOne } from 'typeorm';

import { Channel } from './Channel.js';
import { Integration } from './Integration.js';

@Index('channels_integrations_pkey', ['id'], { unique: true })
@Entity('channels_integrations', { schema: 'public' })
export class ChannelIntegration {
  @Column('text', {
    primary: true,
    name: 'id',
    default: 'gen_random_uuid()',
  })
  id: string;

  @Column('boolean', { name: 'enabled', default: false })
  enabled: boolean;

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

  @Column('jsonb', { name: 'data', nullable: true })
  data: Record<string, any> | null;

  @ManyToOne(() => Channel, (channels) => channels.channelsIntegrations, {
    onDelete: 'RESTRICT',
    onUpdate: 'CASCADE',
  })
  @JoinColumn([{ name: 'channelId', referencedColumnName: 'id' }])
  channel: Channel;

  @ManyToOne(() => Integration, (integrations) => integrations.channelsIntegrations, {
    onDelete: 'CASCADE',
    onUpdate: 'CASCADE',
  })
  @JoinColumn([{ name: 'integrationId', referencedColumnName: 'id' }])
  integration: Integration;
}
