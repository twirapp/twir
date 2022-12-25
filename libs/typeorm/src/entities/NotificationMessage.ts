/* eslint-disable import/no-cycle */
import {
  Column,
  Entity,
  JoinColumn,
  ManyToOne,
  PrimaryColumn,
} from 'typeorm';

import { Notification } from './Notification';

export enum LangCode {
  RU = 'RU',
  GB = 'GB',
}

@Entity('notifications_messages', { schema: 'public' })
export class NotificationMessage {
  @PrimaryColumn('text', {
    primary: true,
    name: 'id',
    default: () => 'gen_random_uuid()',
  })
  id: string;

  @Column('text', { name: 'text' })
  text: string;

  @Column('text', { name: 'title', nullable: true })
  title: string | null;

  @Column('enum', { name: 'langCode', enum: LangCode })
  langCode: LangCode;

  @ManyToOne(() => Notification, _ => _.messages, {
    onDelete: 'RESTRICT',
    onUpdate: 'CASCADE',
  })
  @JoinColumn([{ name: 'notificationId', referencedColumnName: 'id' }])
  notification?: Notification;

  @Column()
  notificationId: string;
}
