import { Column, Entity, JoinColumn, ManyToOne, PrimaryColumn, type Relation } from 'typeorm';

import { type Channel } from './Channel.js';

@Entity('channels_words_counters', { schema: 'public' })
export class ChannelWordCounter {
  @PrimaryColumn('uuid', {
    name: 'id',
    primary: true,
    default: () => 'gen_random_uuid()',
  })
  id: string;

  @Column('text')
  channelId: string;

  @ManyToOne('Channel', 'wordsCounters')
  @JoinColumn({ name: 'channelId' })
  channel: Relation<Channel>;

  @Column('text')
  phrase: string;

  @Column('int4')
  counter: number;

  @Column('bool', { default: true })
  enabled: boolean;
}