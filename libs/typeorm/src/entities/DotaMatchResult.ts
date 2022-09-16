/* eslint-disable import/no-cycle */
import { Column, Entity, Index, JoinColumn, OneToOne } from 'typeorm';

import { DotaMatch } from './DotaMatch.js';

@Index('dota_matches_results_pkey', ['id'], { unique: true })
@Index('dota_matches_results_match_id_key', ['matchId'], { unique: true })
@Entity('dota_matches_results', { schema: 'public' })
export class DotaMatchResult {
  @Column('text', {
    primary: true,
    name: 'id',
    default: 'gen_random_uuid()',
  })
  id: string;

  @Column('text', { name: 'match_id' })
  matchId: string;

  @Column('jsonb', { name: 'players' })
  players: any[];

  @Column('boolean', { name: 'radiant_win' })
  radiantWin: boolean;

  @Column('integer', { name: 'game_mode' })
  gameMode: number;

  @OneToOne(() => DotaMatch, (match) => match.result, {
    onDelete: 'RESTRICT',
    onUpdate: 'CASCADE',
  })
  @JoinColumn([{ name: 'match_id', referencedColumnName: 'matchId' }])
  match: DotaMatch;
}
