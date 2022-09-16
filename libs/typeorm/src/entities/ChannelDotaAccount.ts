/* eslint-disable import/no-cycle */
import { Column, Entity, Index, JoinColumn, ManyToOne } from 'typeorm';

import { Channel } from './Channel.js';

@Index('channels_dota_accounts_id_channelId_key', ['channelId', 'id'], {
  unique: true,
})
@Index('channels_dota_accounts_pkey', ['channelId', 'id'], { unique: true })
@Index('channels_dota_accounts_id_idx', ['id'], {})
@Entity('channels_dota_accounts', { schema: 'public' })
export class DotaAccount {
  @Column('text', { primary: true, name: 'id' })
  id: string;

  @Column('text', { primary: true, name: 'channelId' })
  channelId: string;

  @ManyToOne(() => Channel, (channels) => channels.dotaAccounts, {
    onDelete: 'RESTRICT',
    onUpdate: 'CASCADE',
  })
  @JoinColumn([{ name: 'channelId', referencedColumnName: 'id' }])
  channel: Channel;
}
