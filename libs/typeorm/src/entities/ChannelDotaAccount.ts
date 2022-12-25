/* eslint-disable import/no-cycle */
import {
  Column,
  Entity,
  Index,
  JoinColumn,
  ManyToOne,
  PrimaryColumn,
} from 'typeorm';

import { Channel } from './Channel';

@Index('channels_dota_accounts_id_channelId_key', ['channelId', 'id'], {
  unique: true,
})
@Entity('channels_dota_accounts', { schema: 'public' })
export class ChannelDotaAccount {
  @PrimaryColumn('text', { name: 'id' })
  id: string;

  @ManyToOne(() => Channel, _ => _.dotaAccounts, {
    onDelete: 'RESTRICT',
    onUpdate: 'CASCADE',
  })
  @JoinColumn([{ name: 'channelId', referencedColumnName: 'id' }])
  channel?: Channel;

  @Column('text', { primary: true, name: 'channelId' })
  channelId: string;
}
