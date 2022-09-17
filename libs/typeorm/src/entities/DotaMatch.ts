/* eslint-disable import/no-cycle */
import {
  Column,
  Entity,
  Index,
  JoinColumn,
  ManyToOne,
  OneToMany,
  OneToOne,
  PrimaryColumn,
  // eslint-disable-next-line comma-dangle
  Relation
} from 'typeorm';

import { type DotaGameMode } from './DotaGameMode.js';
import { type DotaMatchCard } from './DotaMatchCard.js';
import { type DotaMatchResult } from './DotaMatchResult.js';

@Entity('dota_matches', { schema: 'public' })
export class DotaMatch {
  @PrimaryColumn('text', {
    primary: true,
    name: 'id',
    default: () => 'gen_random_uuid()',
  })
  id: string;

  @Column('timestamp without time zone', { name: 'startedAt' })
  startedAt: Date;

  @Column('integer', { name: 'lobby_type', nullable: true })
  lobbyType: number | null;

  @Column('int4', { name: 'players', nullable: true, array: true })
  players: number[] | null;

  @Column('int4', { name: 'players_heroes', nullable: true, array: true })
  playersHeroes: number[] | null;

  @Column('text', { name: 'weekend_tourney_bracket_round', nullable: true })
  weekendTourneyBracketRound: string | null;

  @Column('text', { name: 'weekend_tourney_skill_level', nullable: true })
  weekendTourneySkillLevel: string | null;

  @Index()
  @Column('text', { name: 'match_id', unique: true })
  matchId: string;

  @Column('integer', { name: 'avarage_mmr' })
  avarageMmr: number;

  @Column('text', { name: 'lobbyId' })
  lobbyId: string;

  @Column('boolean', { name: 'finished', default: false })
  finished: boolean;

  @ManyToOne('DotaGameMode', 'dotaMatches', {
    onDelete: 'RESTRICT',
    onUpdate: 'CASCADE',
  })
  @JoinColumn([{ name: 'gameModeId', referencedColumnName: 'id' }])
  gameMode?: Relation<DotaGameMode>;

  @Column()
  gameModeId: number;

  @OneToMany('DotaMatchCard', 'match')
  cards?: Relation<DotaMatchCard[]>;

  @OneToOne('DotaMatchResult', 'match')
  result?: Relation<DotaMatchResult>;
}
