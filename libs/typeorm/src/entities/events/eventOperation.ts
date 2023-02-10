import { Column, Entity, JoinColumn, ManyToOne, PrimaryGeneratedColumn } from 'typeorm';

// eslint-disable-next-line import/no-cycle
import { Event } from './event';

export enum OperationType {
  BAN,
  UNBAN,
  BAN_RANDOM,
  VIP,
  UNVIP,
  UNVIP_RANDOM,
  MOD,
  UNMOD,
  UNMOD_RANDOM,
  SEND_MESSAGE,
  CHANGE_TITLE,
  CHANGE_CATEGORY,
  FULFILL_REDEMPTION,
  CANCEL_REDEMPTION,
  ENABLE_SUBMODE,
  DISABLE_SUBMODE,
  ENABLE_EMOTEONLY,
  DISABLE_EMOTEONLY
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
  event: Event;

  @Column('uuid')
  eventId: string;

  @Column('text', { nullable: true })
  input: string | null;

  @Column({ nullable: true })
  repeat: number | null;

  @Column()
  order: number;
}