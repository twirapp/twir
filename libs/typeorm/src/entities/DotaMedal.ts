import { Column, Entity, Index } from 'typeorm';

@Index('dota_medals_rank_tier_idx', ['rankTier'], {})
@Index('dota_medals_pkey', ['rankTier'], { unique: true })
@Entity('dota_medals', { schema: 'public' })
export class DotaMedal {
  @Column('text', { primary: true, name: 'rank_tier' })
  rankTier: string;

  @Column('text', { name: 'name' })
  name: string;
}
