import { Column, Entity, JoinColumn, ManyToOne, PrimaryGeneratedColumn } from 'typeorm';

import { Channel } from './Channel.js';
import { Integration } from './Integration.js';

@Entity('channels_integrations')
export class ChannelIntegration {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @Column({ default: true })
  enabled: boolean;

  @ManyToOne(() => Channel, c => c.integrations)
  @JoinColumn()
  channel: Channel;

  @ManyToOne(() => Integration, i => i.channelIntegrations)
  @JoinColumn()
  integration: Integration;

  @Column()
  accessToken?: string;

  @Column()
  refreshToken?: string;

  @Column()
  apiKey?: string;

  @Column('simple-json')
  data?: Record<string, any>;
}