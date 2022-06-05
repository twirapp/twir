import { Entity, JoinColumn, OneToOne, PrimaryColumn, PrimaryGeneratedColumn } from 'typeorm';

import { User } from './User.js';

@Entity('users_testers')
export class Tester {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @OneToOne(() => User, u => u.tester)
  @JoinColumn()
  user: User;
}