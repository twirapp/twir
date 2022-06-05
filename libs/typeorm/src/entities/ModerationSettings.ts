import { Column, Entity, JoinColumn, OneToOne, PrimaryGeneratedColumn, Unique } from 'typeorm';

import { Channel } from './Channel.js';

export enum SettingsType {
  links = 'links',
  blacklist = 'blacklist',
  symbols = 'symbols',
  longMessages = 'longMessages',
  caps = 'caps',
  emotes = 'emotes'
}

@Entity('channels_moderation_settings')
@Unique('settingsUnique', ['channelId', 'type'])
export class ModerationSettings {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @Column({
    type: 'enum',
    enum: SettingsType,
  })
  type: SettingsType;

  @OneToOne(() => Channel, c => c.moderationSettings)
  @JoinColumn()
  channel: Channel;

  @Column({ default: false })
  enabled: boolean;

  @Column({ default: false })
  subscribers: boolean;

  @Column({ default: false })
  vips: boolean;

  @Column()
  banTime: number;

  @Column()
  banMessage?: string;

  @Column({ default: false })
  checkClips?: boolean;

  @Column()
  triggerLength?: number;

  @Column()
  maxPercentage?: number;

  @Column('simple-array', { default: '[]' })
  blackListSentences?: string[];
}