/* eslint-disable import/no-cycle */
import {
  Column,
  Entity,
  JoinColumn,
  ManyToOne,
  PrimaryGeneratedColumn,
  type Relation,
} from 'typeorm';

import { type Channel } from './Channel.js';
import { type Integration } from './Integration.js';

@Entity('channels_integrations', { schema: 'public' })
export class ChannelIntegration {
  @PrimaryGeneratedColumn('uuid')
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

  @ManyToOne('Channel', 'channelsIntegrations', {
    onDelete: 'RESTRICT',
    onUpdate: 'CASCADE',
  })
  @JoinColumn([{ name: 'channelId', referencedColumnName: 'id' }])
  channel?: Relation<Channel>;

  @Column()
  channelId: string;

  @ManyToOne('Integration', 'channelsIntegrations', {
    onDelete: 'CASCADE',
    onUpdate: 'CASCADE',
  })
  @JoinColumn([{ name: 'integrationId', referencedColumnName: 'id' }])
  integration?: Relation<Integration>;

  @Column()
  integrationId: string;
}
