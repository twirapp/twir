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
  IconPencilPlus,
  IconUser, IconClipboardCopy,
} from '@tabler/icons';
import { AuthUser } from '@tsuwari/shared';
import { useTranslation } from 'next-i18next';
import { useRouter } from 'next/router';
import { useEffect, MouseEvent } from 'react';

import { useProfile } from '@/services/api';
import { createDefaultDashboard, useSelectedDashboard } from '@/services/dashboard';

type Page = {
  label: string;
  icon?: TablerIcon;
  path: string,
  subPages?: Page[]
}

const navigationLinks: Array<Page> = [
  { label: 'Dashboard', icon: IconDashboard, path: '/' },
  { label: 'Integrations', icon: IconBox, path: '/integrations' },
  { label: 'Settings', icon: IconSettings, path: '/settings' },
  {
    label: 'Commands',
    icon: IconCommand,
    path: 'commands',
    subPages: [
      { label: 'Custom', icon: IconPencilPlus, path: '/commands/custom' },
      { label: 'Moderation', icon: IconSword, path: '/commands/moderation' },
      { label: 'Manage', icon: IconClipboardCopy, path: '/commands/manage' },
      { label: 'Dota', path: '/commands/dota' },
    ],
  },
  { label: 'Timers', icon: IconClockHour7, path: '/timers' },
  { label: 'Moderation', icon: IconSword, path: '/moderation' },
  { label: 'Keywords', icon: IconKey, path: '/keywords' },
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

  function goToPageByItem(e: MouseEvent, item: Page) {
    if (item.subPages) return;

    router.push(item.path ? item.path : item.label.toLowerCase());
    props.setOpened(false);
  }

  const computeActive = (item: Page) => {
    if (item.subPages) {
      return router.asPath.startsWith(`/${item.path}`);
    } else {
      return item.path === router.asPath;
    }
  };
  const createNavLink = (item: Page) => <NavLink
    key={item.label}
    active={computeActive(item)}
    label={item.label}
    defaultOpened={item.subPages && router.asPath.startsWith(`/${item.path}`)}
    icon={item.icon ? <item.icon size={16} stroke={1.5} /> : ''}
    onClick={(e) => goToPageByItem(e, item)}
    sx={{ width: '100%' }}
  >{item.subPages && item.subPages.map(p => createNavLink(p))}</NavLink>;

  const links = navigationLinks.map((item) => createNavLink(item));

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
