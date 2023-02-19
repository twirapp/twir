import { Column, Entity, JoinColumn, ManyToOne, PrimaryGeneratedColumn } from 'typeorm';

// eslint-disable-next-line import/no-cycle
import { Event } from './Event';

export enum OperationType {
  TIMEOUT = 'TIMEOUT',
  TIMEOUT_RANDOM = 'TIMEOUT_RANDOM',
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
  DISABLE_EMOTEONLY = 'DISABLE_EMOTEONLY',
  CREATE_GREETING = 'CREATE_GREETING',
  OBS_SET_SCENE = 'OBS_SET_SCENE',
  OBS_TOGGLE_SOURCE = 'OBS_TOGGLE_SOURCE',
  OBS_TOGGLE_AUDIO = 'OBS_TOGGLE_AUDIO',
  OBS_AUDIO_SET_VOLUME = 'OBS_AUDIO_SET_VOLUME',
  OBS_AUDIO_INCREASE_VOLUME = 'OBS_AUDIO_INCREASE_VOLUME',
  OBS_AUDIO_DECREASE_VOLUME = 'OBS_AUDIO_DECREASE_VOLUME',
}

@Entity({ name: 'channels_events_operations' })
export class EventOperation {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @Column('enum', { enum: OperationType })
  type: OperationType;

  @Column({ nullable: false, default: 0 })
  delay: number;

  @ManyToOne(() => Event, _ => _.operations)
  @JoinColumn({ name: 'eventId' })
  event?: Event;

  @Column('uuid')
  eventId: string;

  @Column('text', { nullable: true })
  input: string | null;

  @Column({ default: 1 })
  repeat: number;

  @Column()
  order: number;

  @Column({ default: false })
  useAnnounce: boolean;

  @Column({ default: 600 })
  timeoutTime: number;

  @Column({ nullable: true })
  obsTargetName: string | null;
}