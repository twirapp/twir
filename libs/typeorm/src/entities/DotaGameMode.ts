/* eslint-disable import/no-cycle */
import {
  Column,
  Entity,
  OneToMany,
  PrimaryColumn,
} from 'typeorm';

import { DotaMatch } from './DotaMatch';

@Entity('dota_game_modes', { schema: 'public' })
export class DotaGameMode {
  @PrimaryColumn('integer', { primary: true, name: 'id' })
  id: number;

  @Column('text', { name: 'name' })
  name: string;

  @OneToMany(() => DotaMatch, _ => _.gameMode)
  dotaMatches?: DotaMatch[];
}
