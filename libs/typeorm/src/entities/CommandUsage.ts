import { CreateDateColumn, Entity, JoinColumn, ManyToOne, PrimaryGeneratedColumn } from 'typeorm';

import { Command } from './Command.js';

@Entity('channels_commands_usages')
export class CommandUsage {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @CreateDateColumn()
  createdAt: Date;

  @ManyToOne(() => Command, c => c.usages)
  @JoinColumn()
  command: Command;
}