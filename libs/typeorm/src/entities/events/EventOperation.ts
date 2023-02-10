import { Column, Entity, JoinColumn, ManyToOne, PrimaryGeneratedColumn } from 'typeorm';

// eslint-disable-next-line import/no-cycle
import { Event } from './Event';

export enum OperationType {
  BAN= 'BAN',
  UNBAN = 'UNBAN',
  BAN_RANDOM = 'BAN_RANDOM',
  VIP = 'VIP',
  UNVIP = 'UNVIP',
  UNVIP_RANDOM = 'UNVIP_RANDOM',
  MOD = 'MOD',
  UNMOD = 'UNMOD',
  UNMOD_RANDOM = 'UNMOD_RANDOM',
  SEND_MESSAGE = 'SEND_MESSAGE',
  CHANGE_TITLE = 'CHANGE_TITLE',
  CHANGE_CATEGORY = 'CHANGE_CATEGORY',
  FULFILL_REDEMPTION = 'FULFILL_REDEMPTION',
  CANCEL_REDEMPTION = 'CANCEL_REDEMPTION',
  ENABLE_SUBMODE = 'ENABLE_SUBMODE',
  DISABLE_SUBMODE = 'DISABLE_SUBMODE',
  ENABLE_EMOTEONLY = 'ENABLE_EMOTEONLY',
  DISABLE_EMOTEONLY = 'DISABLE_EMOTEONLY'
}

@Entity({ name: 'channels_events_operations' })
export class EventOperation {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @Column('enum', { enum: OperationType })
  type: OperationType;

  @Column({ nullable: true })
  delay: number | null;

  @ManyToOne(() => Event, _ => _.operations)
  @JoinColumn({ name: 'eventId' })
  event?: Event;

  @Column('uuid')
  eventId: string;

  @Column('text', { nullable: true })
  input: string | null;

  @Column({ nullable: true })
  repeat: number | null;

  @Column()
  order: number;
}