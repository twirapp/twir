/* eslint-disable import/no-cycle */
import { Column, Entity, Index, JoinColumn, ManyToOne, OneToMany } from 'typeorm';

import { NotificationMessage } from './NotificationMessage.js';
import { User } from './User.js';
import { UserViewedNotification } from './UserViewedNotification.js';

@Index('notifications_pkey', ['id'], { unique: true })
@Entity('notifications', { schema: 'public' })
export class Notification {
  @Column('text', {
    primary: true,
    name: 'id',
    default: 'gen_random_uuid()',
  })
  id: string;

  @Column('text', { name: 'imageSrc', nullable: true })
  imageSrc: string | null;

  @Column('timestamp without time zone', {
    name: 'createdAt',
    default: 'CURRENT_TIMESTAMP',
  })
  createdAt: Date;

  @ManyToOne(() => User, (users) => users.notifications, {
    onDelete: 'SET NULL',
    onUpdate: 'CASCADE',
  })
  @JoinColumn([{ name: 'userId', referencedColumnName: 'id' }])
  user: User;

  @OneToMany(() => NotificationMessage, (message) => message.notification)
  messages: NotificationMessage[];

  @OneToMany(
    () => UserViewedNotification,
    (usersViewedNotifications) => usersViewedNotifications.notification,
  )
  viewedNotifications: UserViewedNotification[];
}
