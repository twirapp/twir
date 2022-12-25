/* eslint-disable import/no-cycle */
import {
  Column,
  CreateDateColumn,
  Entity,
  JoinColumn,
  ManyToOne,
  PrimaryColumn,
} from 'typeorm';

import { Notification } from './Notification';
import { User } from './User';

@Entity('users_viewed_notifications', { schema: 'public' })
export class UserViewedNotification {
  @PrimaryColumn('text', {
    primary: true,
    name: 'id',
    default: () => 'gen_random_uuid()',
  })
  id: string;

  @CreateDateColumn()
  createdAt: Date;

  @ManyToOne(() => Notification, _ => _.viewedNotifications, {
    onDelete: 'RESTRICT',
    onUpdate: 'CASCADE',
  })
  @JoinColumn([{ name: 'notificationId', referencedColumnName: 'id' }])
  notification?: Notification;

  @Column()
  notificationId: string;

  @ManyToOne(() => User, _ => _.viewedNotifications, {
    onDelete: 'RESTRICT',
    onUpdate: 'CASCADE',
  })
  @JoinColumn([{ name: 'userId', referencedColumnName: 'id' }])
  user?: User;

  @Column()
  userId: string;
}
