/* eslint-disable import/no-cycle */
import { Column, CreateDateColumn, Entity, Index, JoinColumn, ManyToOne, OneToMany, PrimaryColumn, Relation } from 'typeorm';

import { type NotificationMessage } from './NotificationMessage.js';
import { type User } from './User.js';
import { type UserViewedNotification } from './UserViewedNotification.js';


@Entity('notifications', { schema: 'public' })
export class Notification {
  @PrimaryColumn('text', {
    primary: true,
    name: 'id',
    default: 'gen_random_uuid()',
  })
  id: string;

  @Column('text', { name: 'imageSrc', nullable: true })
  imageSrc: string | null;

  @CreateDateColumn({
    name: 'createdAt',
  })
  createdAt: Date;

  @ManyToOne('User', 'notifications', {
    onDelete: 'SET NULL',
    onUpdate: 'CASCADE',
  })
  @JoinColumn([{ name: 'userId', referencedColumnName: 'id' }])
  user: Relation<User>;

  @OneToMany('NotificationMessage', 'notification')
  messages: Relation<NotificationMessage[]>;

  @OneToMany('UserViewedNotification', 'notification')
  viewedNotifications: Relation<UserViewedNotification[]>;
}
