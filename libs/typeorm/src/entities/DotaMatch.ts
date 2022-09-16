/* eslint-disable import/no-cycle */
import {
  Column,
  Entity,
  Index,
  JoinColumn,
  ManyToOne,
  OneToMany,
  OneToOne,
  Relation,
} from 'typeorm';

import { DotaGameMode } from './DotaGameMode.js';
import { DotaMatchCard } from './DotaMatchCard.js';
import { DotaMatchResult } from './DotaMatchResult.js';

@Index('dota_matches_pkey', ['id'], { unique: true })
@Index('dota_matches_match_id_key', ['matchId'], { unique: true })
@Entity('dota_matches', { schema: 'public' })
export class DotaMatch {
  @Column('text', {
    primary: true,
    name: 'id',
    default: 'gen_random_uuid()',
  })
  id: string;

  @Column('timestamp without time zone', { name: 'startedAt' })
  startedAt: Date;

  @Column('integer', { name: 'lobby_type' })
  lobbyType: number;

  @Column('int4', { name: 'players', nullable: true, array: true })
  players: number[] | null;

  @Column('int4', { name: 'players_heroes', nullable: true, array: true })
  playersHeroes: number[] | null;

  @Column('text', { name: 'weekend_tourney_bracket_round', nullable: true })
  weekendTourneyBracketRound: string | null;

  @Column('text', { name: 'weekend_tourney_skill_level', nullable: true })
  weekendTourneySkillLevel: string | null;

  @Column('text', { name: 'match_id' })
  matchId: string;

  @Column('integer', { name: 'avarage_mmr' })
  avarageMmr: number;

  @Column('text', { name: 'lobbyId' })
  lobbyId: string;

  @Column('boolean', { name: 'finished', default: false })
  finished: boolean;

  @ManyToOne(() => DotaGameMode, (gameMode) => gameMode.dotaMatches, {
    onDelete: 'RESTRICT',
    onUpdate: 'CASCADE',
  })
  @JoinColumn([{ name: 'gameModeId', referencedColumnName: 'id' }])
  gameMode: DotaGameMode;

  @OneToMany(() => DotaMatchCard, (matchCard) => matchCard.match)
  cards: DotaMatchCard[];

  @OneToOne('DotaMatchResult', 'match')
  result: Relation<DotaMatchResult>;
}
