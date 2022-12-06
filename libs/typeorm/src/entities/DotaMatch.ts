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
} from 'typeorm';

import { DotaGameMode } from './DotaGameMode';
import { DotaMatchCard } from './DotaMatchCard';
import { DotaMatchResult } from './DotaMatchResult';

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

  @ManyToOne(() => DotaGameMode, _ => _.dotaMatches, {
    onDelete: 'RESTRICT',
    onUpdate: 'CASCADE',
  })
  @JoinColumn([{ name: 'gameModeId', referencedColumnName: 'id' }])
  gameMode?: DotaGameMode;

  @Column()
  gameModeId: number;

  @OneToMany(() => DotaMatchCard, _ => _.match)
  cards?: DotaMatchCard[];

  @OneToOne(() => DotaMatchResult, _ => _.match)
  result?: DotaMatchResult;
}
