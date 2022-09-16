/* eslint-disable import/no-cycle */
import { Column, Entity, Index, JoinColumn, ManyToOne } from 'typeorm';

import { Notification } from './Notification.js';
import { User } from './User.js';

@Index('users_viewed_notifications_pkey', ['id'], { unique: true })
@Entity('users_viewed_notifications', { schema: 'public' })
export class UserViewedNotification {
  @Column('text', {
    primary: true,
    name: 'id',
    default: 'gen_random_uuid()',
  })
  id: string;

  @Column('timestamp without time zone', {
    name: 'createdAt',
    default: 'CURRENT_TIMESTAMP',
  })
  createdAt: Date;

  @ManyToOne(() => Notification, (notification) => notification.viewedNotifications, {
    onDelete: 'RESTRICT',
    onUpdate: 'CASCADE',
  })
  @JoinColumn([{ name: 'notificationId', referencedColumnName: 'id' }])
  notification: Notification;

  @ManyToOne(() => User, (user) => user.viewedNotifications, {
    onDelete: 'RESTRICT',
    onUpdate: 'CASCADE',
  })
  @JoinColumn([{ name: 'userId', referencedColumnName: 'id' }])
  user: User;
}
