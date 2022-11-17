/* eslint-disable import/no-cycle */
import { Column, Entity, OneToMany, PrimaryColumn, type Relation } from 'typeorm';

import { type DotaMatch } from './DotaMatch.js';

@Entity('dota_game_modes', { schema: 'public' })
export class DotaGameMode {
  @PrimaryColumn('integer', { primary: true, name: 'id' })
  id: number;

  @Column('text', { name: 'name' })
  name: string;

  @OneToMany('DotaMatch', 'gameMode')
  dotaMatches?: Relation<DotaMatch[]>;
}
