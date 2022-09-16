import { Column, Entity, Index } from 'typeorm';

@Index('dota_heroes_id_key', ['id'], { unique: true })
@Index('dota_heroes_pkey', ['id'], { unique: true })
@Entity('dota_heroes', { schema: 'public' })
export class DotaHero {
  @Column('integer', { primary: true, name: 'id' })
  id: number;

  @Column('text', { name: 'name' })
  name: string;
}
