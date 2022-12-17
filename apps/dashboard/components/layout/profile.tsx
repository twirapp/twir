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
  const spotlight = useSpotlight();
  const [selectedDashboard, setSelectedDashboard] = useLocalStorage<Dashboard>({
    key: 'selectedDashboard',
    serialize: (v) => JSON.stringify(v),
    deserialize: (v) => JSON.parse(v),
  });

  function setDefaultDashboard() {
    setSelectedDashboard({
      channelId: props.user.id,
      id: props.user.id,
      twitchUser: props.user,
      userId: props.user.id,
    });
  }

  useEffect(() => {
    spotlight.removeActions(spotlight.actions.map((a) => a.id!));
    const actions = props.user.dashboards
      .filter((d) => d.id != selectedDashboard?.id)
      .map((d) => ({
        title: d.twitchUser.display_name,
        description: d.twitchUser.login,
        onTrigger: () => setSelectedDashboard(d),
        icon: <Avatar src={d.twitchUser.profile_image_url} style={{ borderRadius: 111 }} />,
      }));

    console.log(selectedDashboard);
    if (selectedDashboard?.channelId !== props.user.id) {
      actions.push({
        title: props.user.display_name,
        description: props.user.login,
        onTrigger: () =>
          setSelectedDashboard({
            channelId: props.user.id,
            id: props.user.id,
            twitchUser: props.user,
            userId: props.user.id,
          }),
        icon: <Avatar src={props.user.profile_image_url} style={{ borderRadius: 111 }} />,
      });
    }

    spotlight.registerActions(actions);

    if (!selectedDashboard) {
      setDefaultDashboard();
    } else if (!props.user.dashboards.some((d) => d.id === selectedDashboard.id)) {
      // set default dashboard if user no more have access to selected dashbaord
      setDefaultDashboard();
    } else {
      null;
    }
  }, [props.user]);

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
          <Indicator
            position="bottom-end"
            inline
            label={
              selectedDashboard?.channelId != props.user.id && (
                <Avatar
                  src={props.user.profile_image_url}
                  styles={{
                    image: {
                      width: 25,
                      height: 25,
                      borderRadius: 111,
                      border: '1px solid #77ccdf',
                    },
                  }}
                />
              )
            }
            offset={-5}
            styles={{ indicator: { backgroundColor: 'transparent', marginRight: 10 } }}
          >
            <Avatar
              src={selectedDashboard?.twitchUser.profile_image_url}
              alt={selectedDashboard?.twitchUser.display_name}
              style={{ borderRadius: 111, cursor: 'pointer' }}
            />
          </Indicator>
        </Menu.Target>
        <Menu.Dropdown>
          <Menu.Label>Logged in as {props.user.display_name}</Menu.Label>
          <Menu.Item onClick={() => spotlight.openSpotlight()}>Change channel</Menu.Item>

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
