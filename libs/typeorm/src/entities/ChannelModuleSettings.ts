import {
  Column,
  Entity,
  JoinColumn,
  ManyToOne,
  OneToOne,
  PrimaryGeneratedColumn,
  Relation,
} from 'typeorm';

import { type Channel } from './Channel.js';
import { User } from './User';

export enum ModuleType {
  YOUTUBE_SONG_REQUESTS = 'youtube_song_requests',
  OBS_WEBSOCKET = 'obs_websocket',
  TTS = 'tts',
}

@Entity('channels_modules_settings')
export class ChannelModuleSettings {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @Column('enum', { enum: ModuleType })
  type: ModuleType;

  @Column('jsonb')
  settings: Record<string, any>;

  @Column()
  channelId: string;

  @ManyToOne('Channel', 'modules')
  @JoinColumn({ name: 'channelId' })
  channel?: Relation<Channel>;

  @OneToOne(() => User, _ => _.ttsSettings)
  @JoinColumn({ name: 'userId' })
  user?: User;

  @Column({ nullable: true })
  userId: string | null;
}