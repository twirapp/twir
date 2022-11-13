import { Column, CreateDateColumn, Entity, JoinColumn, ManyToOne, PrimaryColumn, type Relation } from 'typeorm';

import { type Notification } from './Notification.js';
import { type User } from './User.js';

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

  @ManyToOne('Notification', 'viewedNotifications', {
    onDelete: 'RESTRICT',
    onUpdate: 'CASCADE',
  })
  @JoinColumn([{ name: 'notificationId', referencedColumnName: 'id' }])
  notification?: Relation<Notification>;

  @Column()
  notificationId: string;

  @ManyToOne('User', 'viewedNotifications', {
    onDelete: 'RESTRICT',
    onUpdate: 'CASCADE',
  })
  @JoinColumn([{ name: 'userId', referencedColumnName: 'id' }])
  user?: Relation<User>;

  @Column()
  userId: string;
}
