/* eslint-disable import/no-cycle */
import { Column, Entity, JoinColumn, ManyToOne, PrimaryColumn, type Relation } from 'typeorm';

import { type Channel } from './Channel.js';
import { type User } from './User.js';

@Entity('channels_dashboard_access', { schema: 'public' })
export class DashboardAccess {
  @PrimaryColumn('text', {
    primary: true,
    name: 'id',
    default: () => 'gen_random_uuid()',
  })
  id: string;

  @ManyToOne('Channel', 'dashboardAccess', {
    onDelete: 'RESTRICT',
    onUpdate: 'CASCADE',
  })
  @JoinColumn([{ name: 'channelId', referencedColumnName: 'id' }])
  channel?: Relation<Channel>;

  @Column()
  channelId: string;

  @ManyToOne('User', 'dashboardAccess', {
    onDelete: 'RESTRICT',
    onUpdate: 'CASCADE',
  })
  @JoinColumn([{ name: 'userId', referencedColumnName: 'id' }])
  user?: Relation<User>;

  @Column()
  userId: string;
}
