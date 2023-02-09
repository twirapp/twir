import { Column, Entity, ManyToOne, PrimaryGeneratedColumn } from 'typeorm';

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
export class BaseOperation {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @Column('enum', { enum: OperationType })
  type: OperationType;

  @Column({ nullable: true })
  delay: number | null;

  @ManyToOne(() => Event, _ => _.operations)
  event: Event;
}