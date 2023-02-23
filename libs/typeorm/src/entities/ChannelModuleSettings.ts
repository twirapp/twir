import { Column, Entity, JoinColumn, ManyToOne, PrimaryColumn, PrimaryGeneratedColumn, Relation } from 'typeorm';

import { type Channel } from './Channel.js';

export enum ModuleType {
  YOUTUBE_SONG_REQUESTS = 'youtube_song_requests',
  OBS_WEBSOCKET = 'obs_websocket'
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
}