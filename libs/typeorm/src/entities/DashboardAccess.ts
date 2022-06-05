import { Entity, JoinColumn, ManyToOne, PrimaryGeneratedColumn } from 'typeorm';

import { Channel } from './Channel.js';
import { User } from './User.js';

@Entity('channels_dashboard_access')
export class DashboardAccess {
  @PrimaryGeneratedColumn('uuid')

  @ManyToOne(() => Channel, c => c.dashboards)
  @JoinColumn()
  channel: Channel;

  @ManyToOne(() => User, u => u.dashboards)
  @JoinColumn()
  user: User;
}