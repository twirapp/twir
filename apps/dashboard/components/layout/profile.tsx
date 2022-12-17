import { Avatar, Menu, Indicator } from '@mantine/core';
import { useLocalStorage } from '@mantine/hooks';
import { useSpotlight } from '@mantine/spotlight';
import { IconLogout } from '@tabler/icons';
import { AuthUser, Dashboard } from '@tsuwari/shared';
import { useEffect } from 'react';
import useSWRMutation from 'swr/mutation';

import { swrFetcher } from '../../services/swrFetcher';

type Props = {
  user: AuthUser;
};

export function Profile(props: Props) {
  const { trigger: logout } = useSWRMutation('/api/user', async () => {
    const req = await swrFetcher('/api/auth/logout', { method: 'POST' });
    if (req === 'OK') {
      localStorage.removeItem('accessToken');
    }
  });

  return (
    <div>
      <Menu transition="skew-down" shadow="md" width={200}>
        <Menu.Target>
          <Avatar
            src={props.user.profile_image_url}
            alt={props.user.display_name}
            style={{ borderRadius: 111, cursor: 'pointer' }}
          />
        </Menu.Target>
        <Menu.Dropdown>
          <Menu.Label>Logged in as {props.user.display_name}</Menu.Label>

          <Menu.Divider />
          <Menu.Item
            color="red"
            icon={<IconLogout size={14} />}
            onClick={() => logout('/api/auth/logout')}
          >
            Logout
          </Menu.Item>
        </Menu.Dropdown>
      </Menu>
    </div>
  );
}
