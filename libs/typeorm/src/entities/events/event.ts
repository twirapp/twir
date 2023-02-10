import { Column, Entity, OneToMany, PrimaryGeneratedColumn } from 'typeorm';

// eslint-disable-next-line import/no-cycle
import { EventOperation } from './eventOperation';

export enum EventType {
  FOLLOW = 'FOLLOW',
  SUBSCRIBE = 'SUBSCRIBE',
  RESUBSCRIBE = 'RESUBSCRIBE',
  SUB_GIFT = 'SUB_GIFT',
  REDEMPTION_CREATED = 'REDEMPTION_CREATED',
  COMMAND_USED = 'COMMAND_USED',
  FIRST_USER_MESSAGE = 'FIRST_USER_MESSAGE',
  RAIDED = 'RAIDED',
  TITLE_OR_CATEGORY_CHANGED = 'TITLE_OR_CATEGORY_CHANGED',
  STREAM_ONLINE = 'STREAM_ONLINE',
  STREAM_OFFLINE = 'STREAM_OFFLINE',
  ON_CHAT_CLEAR = 'ON_CHAT_CLEAR',
  DONATE = 'DONATE'
}

@Entity({ name: 'channels_events' })
export class Event {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @Column('enum', { enum: EventType })
  type: EventType;

  @Column('uuid', { nullable: true })
  rewardId: string | null;

  @Column('text', { nullable: true })
  commandId: string | null;

  @OneToMany(() => EventOperation, _ => _.event)
  operations: EventOperation[];
}