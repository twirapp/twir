import { Anchor, Avatar, Flex, Group, Tabs, TabsValue, Text } from '@mantine/core';
import { useRouter } from 'next/router';
import { useEffect, useState } from 'react';

import { useUsersByNames } from '@/services/users';

const routes = ['commands', 'song-requests'] as const;

export const Header: React.FC = () => {
  const router = useRouter();
  const { data: users } = useUsersByNames([router.query.channelName as string]);
  const [activeTab, setActiveTab] = useState<string | null>('commands');

  useEffect(() => {
    const route = routes.find(r => router.route.includes(r));
    if (route) {
      setActiveTab(route);
    }
  }, [router]);

  async function onTabChange(t: TabsValue) {
    const { channelName } = router.query;
    if (!t || !channelName) return;
    await router.push(`/${channelName}/${t}`);
    setActiveTab(t);
  }

  return (
    <Flex gap={'lg'} justify={'space-between'}>
      <Group>
      <Avatar
        size={'xl'}
        src={users?.at(0)?.profile_image_url}
        radius={111}
        alt={users?.at(0)?.display_name}
      />
      <Flex direction={'column'}>
        <Text size={'xl'}>{users?.at(0) && users.at(0)!.display_name}</Text>
        <Anchor>{users?.at(0) && `twitch.tv/${users.at(0)!.login}`}</Anchor>
      </Flex>

      </Group>

      <Tabs
        value={activeTab}
        onTabChange={onTabChange}
      >
        <Tabs.List>
          <Tabs.Tab value="commands">Commands</Tabs.Tab>
          <Tabs.Tab value="song-requests">Song requests</Tabs.Tab>
        </Tabs.List>
      </Tabs>
    </Flex>
  );
};