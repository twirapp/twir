import { Avatar, Menu } from '@mantine/core';
import { IconLogout } from '@tabler/icons';
import { type Profile as IProfile } from '@twir/grpc/generated/api/api/auth';

import { useLogoutMutation } from '@/services/api';

type Props = {
  user: IProfile;
};

export function Profile(props: Props) {
  const logout = useLogoutMutation();

  return (
    <div>
      <Menu transition="pop" shadow="md" withArrow width={200}>
        <Menu.Target>
          <Avatar
            size={34}
            radius="xs"
            style={{ cursor: 'pointer' }}
            src={props.user.avatar}
            alt={props.user.displayName}
          />
        </Menu.Target>
        <Menu.Dropdown>
          <Menu.Label>Logged in as {props.user.displayName}</Menu.Label>
          <Menu.Divider />
          <Menu.Item color="red" icon={<IconLogout size={14} />} onClick={() => logout.mutate()}>
            Logout
          </Menu.Item>
        </Menu.Dropdown>
      </Menu>
    </div>
  );
}
