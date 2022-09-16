/* eslint-disable import/no-cycle */
import { Column, Entity, Index, OneToMany } from 'typeorm';

import { DotaMatch } from './DotaMatch.js';

@Index('dota_game_modes_pkey', ['id'], { unique: true })
@Entity('dota_game_modes', { schema: 'public' })
export class DotaGameMode {
  @Column('integer', { primary: true, name: 'id' })
  id: number;

  @Column('text', { name: 'name' })
  name: string;

  @OneToMany(() => DotaMatch, (match) => match.gameMode)
  dotaMatches: DotaMatch[];
}
