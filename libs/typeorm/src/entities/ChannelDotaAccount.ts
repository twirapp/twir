/* eslint-disable import/no-cycle */
import { Column, Entity, Index, JoinColumn, ManyToOne, PrimaryColumn, Relation } from 'typeorm';

import { type Channel } from './Channel.js';

@Index('channels_dota_accounts_id_channelId_key', ['channelId', 'id'], {
  unique: true,
})
@Entity('channels_dota_accounts', { schema: 'public' })
export class DotaAccount {
  @PrimaryColumn('text', { name: 'id' })
  id: string;

  @Column('text', { primary: true, name: 'channelId' })
  channelId: string;

  @ManyToOne('Channel', 'dotaAccounts', {
    onDelete: 'RESTRICT',
    onUpdate: 'CASCADE',
  })
  @JoinColumn([{ name: 'channelId', referencedColumnName: 'id' }])
  channel: Relation<Channel>;
}
