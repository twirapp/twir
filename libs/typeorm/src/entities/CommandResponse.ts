import { Column, Entity, ManyToOne, PrimaryGeneratedColumn } from 'typeorm';

import { Command } from './Command.js';

@Entity('channels_commands_responses')
export class CommandResponse {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @Column()
  text: string;

  @ManyToOne(() => Command, c => c.responses)
  command: Command;
}