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
import { useViewportSize } from '@mantine/hooks';
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
import { AuthUser } from '@tsuwari/shared';
import { useTranslation } from 'next-i18next';
import { useRouter } from 'next/router';
import { useEffect } from 'react';

import { useProfile } from '@/services/api';
import { createDefaultDashboard, useSelectedDashboard } from '@/services/dashboard';

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

type Props = {
  opened: boolean,
  setOpened: React.Dispatch<React.SetStateAction<boolean>>
}

export function SideBar(props: Props) {
  const viewPort = useViewportSize();
  const router = useRouter();
  const theme = useMantineTheme();
  const { t } = useTranslation('layout');

  const { data: user } = useProfile();

  const spotlight = useSpotlight();
  const [selectedDashboard, setSelectedDashboard] = useSelectedDashboard();

  const setDefaultDashboard = (user: AuthUser) =>
    setSelectedDashboard(createDefaultDashboard(user));

  function openSpotlight() {
    if (user && selectedDashboard) {
      spotlight.removeActions(spotlight.actions.map((a) => a.id!));
      const actions = user.dashboards
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

      if (selectedDashboard.channelId != user.id) {
        actions.unshift({
          title: user.display_name,
          description: user.login,
          onTrigger: () => {
            setSelectedDashboard({
              channelId: user.id,
              id: user.id,
              twitchUser: user,
              userId: user.id,
            });
          },
          icon: <Avatar src={user.profile_image_url} style={{ borderRadius: 111 }} />,
          id: user.id,
        });
      }

      spotlight.registerActions(actions);
      spotlight.openSpotlight();
    }
  }

  useEffect(() => {
    if (!user) return;

    if (!selectedDashboard) {
      return setDefaultDashboard(user);
    }
  }, [user]);

  const links = navigationLinks.map((item) => (
    <NavLink
      key={item.label}
      active={item.path === router.asPath}
      label={item.label}
      icon={<item.icon size={16} stroke={1.5} />}
      onClick={(e) => {
        e.preventDefault();
        router.push(item.path ? item.path : item.label.toLowerCase());
        props.setOpened(false);
      }}
      sx={{ width: '100%' }}
    />
  ));

  return (
    <Navbar hiddenBreakpoint="sm" hidden={!props.opened} width={{ sm: 150, lg: 250 }}>
      <Navbar.Section grow>
        <ScrollArea.Autosize
          maxHeight={viewPort.height - 120}
          type="auto"
          offsetScrollbars={true}
          styles={{
            viewport: {
              padding: 0,
            },
          }}
        >
          <Box component={ScrollArea} sx={{ width: '100%' }}>
            {links}
          </Box>
        </ScrollArea.Autosize>
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
                  {t('sidebar.manage')}
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
