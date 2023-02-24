/* eslint-disable import/no-cycle */
import { Column, Entity, JoinColumn, ManyToOne, PrimaryGeneratedColumn } from 'typeorm';

import { ChannelRole } from './ChannelRole';
import { User } from './User';

@Entity('channels_roles_users')
export class ChannelRoleUser {
    @PrimaryGeneratedColumn('uuid')
    id: number;

    @Column()
    userId: number;

    @Column()
    roleId: number;

    @ManyToOne(() => ChannelRole, _ => _.users, { onDelete: 'CASCADE' })
    @JoinColumn({ name: 'roleId' })
    role?: ChannelRole;

    @ManyToOne(() => User, _ => _.channelRoleUsers)
    @JoinColumn({ name: 'userId' })
    user?: User;
}