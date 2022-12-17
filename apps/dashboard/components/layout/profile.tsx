import { Avatar, Image, Menu, ScrollArea, TextInput, Loader } from '@mantine/core';
import { IconLogout } from '@tabler/icons';
import { AuthUser } from '@tsuwari/shared';
import { useState } from 'react';

type Props = {
  user: AuthUser;
};

export function Profile(props: Props) {
  const [searchState, setSearch] = useState('');

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
          <Menu.Label>S4tont</Menu.Label>
          <TextInput
            placeholder="Search user..."
            value={searchState}
            onChange={(event) => setSearch(event.currentTarget.value)}
            style={{ marginBottom: 5 }}
          />
          <ScrollArea type="auto" style={{ height: 250 }}>
            {props.user.dashboards
              .filter((d) => (searchState !== '' ? d.twitchUser.login.includes(searchState) : true))
              .map((d) => (
                <Menu.Item
                  icon={<Image src={d.twitchUser.profile_image_url} height={20} />}
                  key={d.userId + d.id}
                >
                  {d.twitchUser.login}
                </Menu.Item>
              ))}
          </ScrollArea>

          <Menu.Divider />
          <Menu.Item color="red" icon={<IconLogout size={14} />}>
            Logout
          </Menu.Item>
        </Menu.Dropdown>
      </Menu>
    </div>
  );
}
