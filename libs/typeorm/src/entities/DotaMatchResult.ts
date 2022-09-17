/* eslint-disable import/no-cycle */
import { Column, Entity, Index, JoinColumn, OneToOne, PrimaryColumn, Relation } from 'typeorm';

import { type DotaMatch } from './DotaMatch.js';

@Index('dota_matches_results_match_id_key', ['matchId'], { unique: true })
@Entity('dota_matches_results', { schema: 'public' })
export class DotaMatchResult {
  @PrimaryColumn('text', {
    primary: true,
    name: 'id',
    default: () => 'gen_random_uuid()',
  })
  id: string;

  @Column('jsonb', { name: 'players' })
  players: any[];

  @Column('boolean', { name: 'radiant_win' })
  radiantWin: boolean;

  @Column('integer', { name: 'game_mode' })
  gameMode: number;

  @OneToOne('DotaMatch', 'result', {
    onDelete: 'RESTRICT',
    onUpdate: 'CASCADE',
  })
  @JoinColumn([{ name: 'match_id', referencedColumnName: 'matchId' }])
  match?: Relation<DotaMatch>;

  @Column('text', { name: 'match_id' })
  matchId: string;
}
