import {
  Column,
  Entity,
  Index,
  PrimaryColumn,
} from 'typeorm';

@Entity('dota_heroes', { schema: 'public' })
export class DotaHero {
  @PrimaryColumn('integer', { primary: true, name: 'id' })
  id: number;

  @Index()
  @Column('text', { name: 'name' })
  name: string;
}
