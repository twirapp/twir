/* eslint-disable import/no-cycle */
import { Column, Entity, Index, JoinColumn, ManyToOne } from 'typeorm';

import { Notification } from './Notification.js';

export enum LangCode {
  RU = 'RU',
  GB = 'GB',
}

@Index('notifications_messages_pkey', ['id'], { unique: true })
@Entity('notifications_messages', { schema: 'public' })
export class NotificationMessage {
  @Column('text', {
    primary: true,
    name: 'id',
    default: 'gen_random_uuid()',
  })
  id: string;

  @Column('text', { name: 'text' })
  text: string;

  @Column('text', { name: 'title', nullable: true })
  title: string | null;

  @Column('enum', { name: 'langCode', enum: LangCode })
  langCode: LangCode;

  @ManyToOne(() => Notification, (notifications) => notifications.messages, {
    onDelete: 'RESTRICT',
    onUpdate: 'CASCADE',
  })
  @JoinColumn([{ name: 'notificationId', referencedColumnName: 'id' }])
  notification: Notification;
}
