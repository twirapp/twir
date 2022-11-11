import { Column, Entity, JoinColumn, ManyToOne, PrimaryColumn, Relation } from 'typeorm';

import { type Channel } from './Channel.js';

export enum ModuleType {
  YOUTUBE_SONG_REQUESTS = 'youtube_song_requests',
}

@Entity('channels_modules_settings')
export class ChannelModuleSettings {
  @PrimaryColumn('uuid', {
    primary: true,
    name: 'id',
    default: () => 'gen_random_uuid()',
  })
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

export type YoutubeSettings = {
  maxRequests?: number;
  acceptOnlyWhenOnline?: boolean;
  channelPointsRewardName?: string;
  user?: {
    maxRequests?: number;
    minWatchTime?: number;
    minMessages?: number;
    minFollowTime?: number;
  };
  song?: {
    maxLength?: number;
    minViews?: number;
    acceptedCategories?: string[];
  };
  blackList?: {
    usersIds?: string[]; // it will be converted via twitch api. We'll get user by name and store id
    songsIds?: string[];
    channelsIds?: string[];
    artistsNames?: string[];
  };
};
