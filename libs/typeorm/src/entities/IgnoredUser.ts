import { Column, Entity, Index, PrimaryColumn } from 'typeorm';

@Entity('users_ignored', { schema: 'public' })
export class IgnoredUser {
  @PrimaryColumn('text', { primary: true, name: 'id' })
  id: string;

  @Index()
  @Column('text', { name: 'login', nullable: true })
  login: string | null;

  @Index()
  @Column('text', { name: 'displayName', nullable: true })
  displayName: string | null;
}
