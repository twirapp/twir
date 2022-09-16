/* eslint-disable import/no-cycle */
import { Column, CreateDateColumn, Entity, Index, JoinColumn, ManyToOne, PrimaryColumn, Relation } from 'typeorm';

import { type Notification } from './Notification.js';
import { type User } from './User.js';


@Entity('users_viewed_notifications', { schema: 'public' })
export class UserViewedNotification {
  @PrimaryColumn('text', {
    primary: true,
    name: 'id',
    default: 'gen_random_uuid()',
  })
  id: string;

  @CreateDateColumn()
  createdAt: Date;

  @ManyToOne('Notification', 'viewedNotifications', {
    onDelete: 'RESTRICT',
    onUpdate: 'CASCADE',
  })
  @JoinColumn([{ name: 'notificationId', referencedColumnName: 'id' }])
  notification: Relation<Notification>;

  @ManyToOne('User', 'viewedNotifications', {
    onDelete: 'RESTRICT',
    onUpdate: 'CASCADE',
  })
  @JoinColumn([{ name: 'userId', referencedColumnName: 'id' }])
  user: Relation<User>;
}
