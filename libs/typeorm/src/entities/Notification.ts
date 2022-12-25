/* eslint-disable import/no-cycle */
import {
  Column,
  CreateDateColumn,
  Entity,
  JoinColumn,
  ManyToOne,
  OneToMany,
  PrimaryColumn,
} from 'typeorm';

import { NotificationMessage } from './NotificationMessage';
import { User } from './User';
import { UserViewedNotification } from './UserViewedNotification';

@Entity('notifications', { schema: 'public' })
export class Notification {
  @PrimaryColumn('text', {
    primary: true,
    name: 'id',
    default: () => 'gen_random_uuid()',
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
  user?: User;

  @Column({ nullable: true })
  userId?: string | null;

  @OneToMany(() => NotificationMessage, _ => _.notification)
  messages?: NotificationMessage[];

  @OneToMany(() => UserViewedNotification, _ => _.notification)
  viewedNotifications?: UserViewedNotification[];
}
