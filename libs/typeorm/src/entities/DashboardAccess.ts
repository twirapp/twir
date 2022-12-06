/* eslint-disable import/no-cycle */
import {
  Column,
  Entity,
  JoinColumn,
  ManyToOne,
  PrimaryColumn,
} from 'typeorm';

import { Channel } from './Channel';
import { User } from './User';

@Entity('channels_dashboard_access', { schema: 'public' })
export class DashboardAccess {
  @PrimaryColumn('text', {
    primary: true,
    name: 'id',
    default: () => 'gen_random_uuid()',
  })
  id: string;

  @ManyToOne(() => Channel, _ => _.dashboardAccess, {
    onDelete: 'RESTRICT',
    onUpdate: 'CASCADE',
  })
  @JoinColumn([{ name: 'channelId', referencedColumnName: 'id' }])
  channel?: Channel;

  @Column()
  channelId: string;

  @ManyToOne(() => User, _ => _.dashboardAccess, {
    onDelete: 'RESTRICT',
    onUpdate: 'CASCADE',
  })
  @JoinColumn([{ name: 'userId', referencedColumnName: 'id' }])
  user?: User;

  @Column()
  userId: string;
}
