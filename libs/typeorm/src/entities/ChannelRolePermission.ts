import { Column, Entity, JoinColumn, ManyToOne, PrimaryGeneratedColumn } from 'typeorm';

// eslint-disable-next-line import/no-cycle
import { ChannelRole } from './ChannelRole';
// eslint-disable-next-line import/no-cycle
import { RolePermission } from './RolePermission';

@Entity('channels_roles_permissions')
export class ChannelRolePermission {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @Column()
  roleId: string;

  @Column()
  permissionId: string;

  @ManyToOne(() => ChannelRole, _ => _.permissions, { onDelete: 'CASCADE' })
  @JoinColumn({ name: 'roleId' })
  role?: ChannelRole;

  @ManyToOne(() => RolePermission, _ => _.channelRolePermissions)
  @JoinColumn({ name: 'permissionId' })
  permission?: RolePermission;
}