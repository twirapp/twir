/* eslint-disable import/no-cycle */
import { Column, Entity, Index, JoinColumn, ManyToOne } from 'typeorm';

import { Channel } from './Channel.js';
import { User } from './User.js';

@Index('channels_dashboard_access_pkey', ['id'], { unique: true })
@Entity('channels_dashboard_access', { schema: 'public' })
export class DashboardAccess {
  @Column('text', {
    primary: true,
    name: 'id',
    default: 'gen_random_uuid()',
  })
  id: string;

  @ManyToOne(() => Channel, (channels) => channels.dashboardAccess, {
    onDelete: 'RESTRICT',
    onUpdate: 'CASCADE',
  })
  @JoinColumn([{ name: 'channelId', referencedColumnName: 'id' }])
  channel: Channel;

  @ManyToOne(() => User, (users) => users.dashboardAccess, {
    onDelete: 'RESTRICT',
    onUpdate: 'CASCADE',
  })
  @JoinColumn([{ name: 'userId', referencedColumnName: 'id' }])
  user: User;
}
