/* eslint-disable import/no-cycle */
import { Column, Entity, Index, JoinColumn, ManyToOne, PrimaryColumn, Relation } from 'typeorm';

import { type DotaMatch } from './DotaMatch.js';

@Index('dota_matches_cards_match_id_account_id_key', ['accountId', 'matchId'], {
  unique: true,
})
@Entity('dota_matches_cards', { schema: 'public' })
export class DotaMatchCard {
  @PrimaryColumn('text', {
    primary: true,
    name: 'id',
    default: () => 'gen_random_uuid()',
  })
  id: string;

  @Column('text', { name: 'account_id' })
  accountId: string;

  @Column('integer', { name: 'rank_tier', nullable: true })
  rankTier: number | null;

  @Column('integer', { name: 'leaderboard_rank', nullable: true })
  leaderboardRank: number | null;

  @ManyToOne('DotaMatch', 'cards', {
    onDelete: 'RESTRICT',
    onUpdate: 'CASCADE',
  })
  @JoinColumn([{ name: 'match_id', referencedColumnName: 'id' }])
  match?: Relation<DotaMatch>;

  @Column('text', { name: 'match_id' })
  matchId: string;
}
