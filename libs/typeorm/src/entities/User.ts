import { Entity, JoinColumn, OneToMany, OneToOne, PrimaryColumn } from 'typeorm';

import { Channel } from './Channel.js';
import { DashboardAccess } from './DashboardAccess.js';
import { Greeting } from './Greeting.js';
import { Permit } from './Permit.js';
import { Tester } from './Tester.js';
import { TwitchToken } from './TwitchToken.js';
import { UserStats } from './UserStats.js';

@Entity('users')
export class User {
  @PrimaryColumn()
  id: string;

  @OneToOne(() => Channel, (c) => c.user)
  channel: Channel;

  @OneToOne(() => UserStats, (s) => s.user)
  stats: UserStats;

  @OneToOne(() => Tester, t => t.user)
  tester: Tester;

  @OneToMany(() => DashboardAccess, d => d.user)
  dashboards: DashboardAccess[];

  @OneToOne(() => TwitchToken, t => t.user)
  @JoinColumn()
  token?: TwitchToken;

  @OneToMany(() => Greeting, t => t.user)
  greetings: Greeting[];

  @OneToMany(() => Permit, p => p.user)
  permits: Permit[];
}