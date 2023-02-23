import { Column, Entity, JoinColumn, ManyToOne, PrimaryGeneratedColumn } from 'typeorm';

// eslint-disable-next-line import/no-cycle
import { ChannelRole } from './ChannelRole';
// eslint-disable-next-line import/no-cycle
import { RoleFlag } from './RoleFlag';

@Entity('channels_roles_permissions')
export class ChannelRolePermission {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @Column()
  roleId: string;

  @Column()
  flagId: string;

  @ManyToOne(() => ChannelRole, _ => _.permissions, { onDelete: 'CASCADE' })
  @JoinColumn({ name: 'roleId' })
  role?: ChannelRole;

  @ManyToOne(() => RoleFlag, _ => _.channelRolePermissions)
  @JoinColumn({ name: 'flagId' })
  flag?: RoleFlag;
}