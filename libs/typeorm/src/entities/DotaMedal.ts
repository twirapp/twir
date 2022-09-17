import { Column, Entity, Index, PrimaryColumn } from 'typeorm';

@Entity('dota_medals', { schema: 'public' })
export class DotaMedal {
  @PrimaryColumn('text', { primary: true, name: 'rank_tier' })
  rankTier: string;

  @Index()
  @Column('text', { name: 'name' })
  name: string;
}
