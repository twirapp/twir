import {
  Box,
  Group,
  Navbar,
  NavLink,
  ScrollArea,
  UnstyledButton,
  useMantineTheme,
  Text,
  Avatar,
} from '@mantine/core';
import { useLocalStorage } from '@mantine/hooks';
import { useSpotlight } from '@mantine/spotlight';
import {
  IconDashboard,
  IconBox,
  IconSettings,
  IconCommand,
  IconClockHour7,
  IconSword,
  IconKey,
  IconActivity,
  IconSpeakerphone,
  TablerIcon,
} from '@tabler/icons';
import { AuthUser, Dashboard } from '@tsuwari/shared';
import { useRouter } from 'next/router';
import { useEffect } from 'react';
import useSWR from 'swr';

import { swrFetcher } from '../../services/swrFetcher';

const navigationLinks: Array<{ label: string; icon: TablerIcon; path: string }> = [
  { label: 'Dashboard', icon: IconDashboard, path: '/' },
  { label: 'Integrations', icon: IconBox, path: '/integrations' },
  { label: 'Settings', icon: IconSettings, path: '/settings' },
  { label: 'Commands', icon: IconCommand, path: '/commands' },
  { label: 'Timers', icon: IconClockHour7, path: '/timers' },
  { label: 'Moderation', icon: IconSword, path: '/moderation' },
  { label: 'Keywords', icon: IconKey, path: 'keywords' },
  { label: 'Variables', icon: IconActivity, path: '/variables' },
  { label: 'Greetings', icon: IconSpeakerphone, path: '/greetings' },
];

export function SideBar({ opened }: { opened: boolean }) {
  const router = useRouter();
  const theme = useMantineTheme();

  const { data: userData } = useSWR<AuthUser>('/api/auth/profile', swrFetcher);

  const spotlight = useSpotlight();
  const [selectedDashboard, setSelectedDashboard] = useLocalStorage<Dashboard>({
    key: 'selectedDashboard',
    serialize: (v) => JSON.stringify(v),
    deserialize: (v) => JSON.parse(v),
  });

  function setDefaultDashboard() {
    setSelectedDashboard({
      channelId: userData!.id,
      id: userData!.id,
      twitchUser: userData!,
      userId: userData!.id,
    });
  }

  function openSpotlight() {
    if (userData) {
      spotlight.removeActions(spotlight.actions.map((a) => a.id!));
      const actions = userData.dashboards
        .filter((d) => d.id != selectedDashboard?.id)
        .map((d) => ({
          title: d.twitchUser.display_name,
          description: d.twitchUser.login,
          onTrigger: () => {
            setSelectedDashboard(d);
          },
          icon: <Avatar src={d.twitchUser.profile_image_url} style={{ borderRadius: 111 }} />,
          id: d.id,
        }));

      if (selectedDashboard.channelId != userData.id) {
        actions.unshift({
          title: userData.display_name,
          description: userData.login,
          onTrigger: () => {
            setSelectedDashboard({
              channelId: userData.id,
              id: userData.id,
              twitchUser: userData,
              userId: userData.id,
            });
          },
          icon: <Avatar src={userData.profile_image_url} style={{ borderRadius: 111 }} />,
          id: userData.id,
        });
      }

      spotlight.registerActions(actions);
      spotlight.openSpotlight();
    }
  }

  useEffect(() => {
    if (userData) {
      if (!selectedDashboard) {
        setDefaultDashboard();
      } else if (!userData.dashboards.some((d) => d.id === selectedDashboard.id)) {
        // set default dashboard if user no more have access to selected dashbaord
        setDefaultDashboard();
      } else {
        null;
      }
    }
  }, [userData]);

  const links = navigationLinks.map((item, index) => (
    <NavLink
      key={item.label}
      active={item.path === router.asPath}
      label={item.label}
      icon={<item.icon size={16} stroke={1.5} />}
      onClick={(e) => {
        e.preventDefault();
        router.push(item.path ? item.path : item.label.toLowerCase());
      }}
    />
  ));

  return (
    <Navbar hiddenBreakpoint="sm" hidden={!opened} width={{ sm: 150, lg: 250 }}>
      <Navbar.Section grow>
        <Box component={ScrollArea}>{links}</Box>
      </Navbar.Section>
      <Navbar.Section>
        <Box
          sx={{
            padding: theme.spacing.sm,
            borderTop: `1px solid ${
              theme.colorScheme === 'dark' ? theme.colors.dark[4] : theme.colors.gray[2]
            }`,
          }}
        >
          <UnstyledButton
            sx={{
              display: 'block',
              width: '100%',
              padding: theme.spacing.xs,
              borderRadius: theme.radius.sm,
              color: theme.colorScheme === 'dark' ? theme.colors.dark[0] : theme.black,

              '&:hover': {
                backgroundColor:
                  theme.colorScheme === 'dark' ? theme.colors.dark[6] : theme.colors.gray[0],
              },
            }}
            onClick={openSpotlight}
          >
            <Group>
              <Avatar src={selectedDashboard?.twitchUser.profile_image_url} radius="xl" />
              <Box sx={{ flex: 1 }}>
                <Text size="xs" weight={500}>
                  Managing channel
                </Text>
                <Text color="dimmed" size="xs">
                  {selectedDashboard?.twitchUser.display_name}
                </Text>
              </Box>
            </Group>
          </UnstyledButton>
        </Box>
      </Navbar.Section>
    </Navbar>
  );
}
