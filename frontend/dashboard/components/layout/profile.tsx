import { Avatar, Menu } from '@mantine/core';
import { IconLogout } from '@tabler/icons';
import { AuthUser } from '@tsuwari/shared';

import { useLogoutMutation } from '@/services/api';

type Props = {
  user: AuthUser;
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
            src={props.user.profile_image_url}
            alt={props.user.display_name}
          />
        </Menu.Target>
        <Menu.Dropdown>
          <Menu.Label>Logged in as {props.user.display_name}</Menu.Label>
          <Menu.Divider />
          <Menu.Item color="red" icon={<IconLogout size={14} />} onClick={() => logout.mutate()}>
            Logout
          </Menu.Item>
        </Menu.Dropdown>
      </Menu>
    </div>
  );
}
