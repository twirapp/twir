import { Column, Entity, Index, JoinColumn, ManyToOne } from 'typeorm';

// eslint-disable-next-line import/no-cycle
import { User } from './User.js';

@Index('users_files_pkey', ['id'], { unique: true })
@Entity('users_files', { schema: 'public' })
export class UserFile {
  @Column('text', {
    primary: true,
    name: 'id',
    default: 'gen_random_uuid()',
  })
  id: string;

  @Column('text', { name: 'name' })
  name: string;

  @Column('integer', { name: 'size' })
  size: number;

  @Column('text', { name: 'type' })
  type: string;

  @ManyToOne(() => User, (user) => user.files, {
    onDelete: 'SET NULL',
    onUpdate: 'CASCADE',
  })
  @JoinColumn([{ name: 'userId', referencedColumnName: 'id' }])
  user: User;
}
