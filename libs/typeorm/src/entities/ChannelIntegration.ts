/* eslint-disable import/no-cycle */
import {
  Column,
  Entity,
  JoinColumn,
  ManyToOne,
  PrimaryGeneratedColumn,
} from 'typeorm';

import { Channel } from './Channel';
import { Integration } from './Integration';

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

  @ManyToOne(() => Channel, _ => _.channelsIntegrations, {
    onDelete: 'RESTRICT',
    onUpdate: 'CASCADE',
  })
  @JoinColumn([{ name: 'channelId', referencedColumnName: 'id' }])
  channel?: Channel;

  @Column()
  channelId: string;

  @ManyToOne(() => Integration, _ => _.channelsIntegrations, {
    onDelete: 'CASCADE',
    onUpdate: 'CASCADE',
  })
  @JoinColumn([{ name: 'integrationId', referencedColumnName: 'id' }])
  integration?: Integration;

  @Column()
  integrationId: string;
}
