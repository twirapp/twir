import { Column, Entity, JoinColumn, ManyToOne, OneToMany, OneToOne, PrimaryColumn } from 'typeorm';

import { Bot } from './Bot.js';
import { ChannelIntegration } from './ChannelIntegration.js';
import { Command } from './Command.js';
import { DashboardAccess } from './DashboardAccess.js';
import { Greeting } from './Greeting.js';
import { Keyword } from './Keyword.js';
import { ModerationSettings } from './ModerationSettings.js';
import { Permit } from './Permit.js';
import { Timer } from './Timer.js';
import { User } from './User.js';
import { UserStats } from './UserStats.js';

@Entity('channels')
export class Channel {
  @PrimaryColumn()
  id: string;

  @Column({ default: true })
  isEnabled: boolean;

  @Column({ default: false })
  isTwitchBanned: boolean;

  @Column({ default: false })
  isBanned: boolean;

  @OneToOne(() => User, (u) => u.channel)
  @JoinColumn()
  user: User;

  @OneToOne(() => UserStats, (s) => s.channel)
  usersStats: UserStats;

  @OneToMany(() => Command, (c) => c.channel)
  commands: Command[];

  @OneToMany(() => DashboardAccess, d => d.channel)
  dashboards: DashboardAccess[];

  @OneToMany(() => ChannelIntegration, i => i.channel)
  integrations: ChannelIntegration[];

  @ManyToOne(() => Bot, b => b.channels)
  @JoinColumn()
  bot: Bot;

  @OneToMany(() => Timer, t => t.channel)
  timers: Timer[];

  @OneToMany(() => Greeting, t => t.channel)
  greetings: Greeting[];

  @OneToOne(() => ModerationSettings, s => s.channel)
  moderationSettings: ModerationSettings;

  @OneToMany(() => Permit, p => p.channel)
  permits: Permit[];

  @OneToMany(() => Keyword, k => k.channel)
  keywords: Keyword[];
}