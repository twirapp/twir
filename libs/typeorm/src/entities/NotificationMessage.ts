/* eslint-disable import/no-cycle */
import { Column, Entity, Index, JoinColumn, ManyToOne, PrimaryColumn, Relation } from 'typeorm';

import { type Notification } from './Notification.js';

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

  @ManyToOne('Notification', 'messages', {
    onDelete: 'RESTRICT',
    onUpdate: 'CASCADE',
  })
  @JoinColumn([{ name: 'notificationId', referencedColumnName: 'id' }])
  notification?: Relation<Notification>;

  @Column()
  notificationId: string;
}
