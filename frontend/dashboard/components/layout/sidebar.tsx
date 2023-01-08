import {
  Avatar,
  Box,
  createStyles,
  Group,
  Navbar,
  NavLink,
  ScrollArea,
  Text,
  UnstyledButton,
  useMantineTheme,
} from '@mantine/core';
import { useViewportSize } from '@mantine/hooks';
import { useSpotlight } from '@mantine/spotlight';
import {
  IconActivity,
  IconBox,
  IconClipboardCopy,
  IconClockHour7,
  IconCommand,
  IconDashboard,
  IconHeadphones,
  IconKey,
  IconPencilPlus,
  IconPlayerPlay,
  IconPlaylist,
  IconSettings,
  IconSpeakerphone,
  IconSword,
  TablerIcon,
} from '@tabler/icons';
import { Dashboard } from '@tsuwari/shared';
import { useTranslation } from 'next-i18next';
import { useRouter } from 'next/router';
import { useCallback, useContext, useEffect, useState } from 'react';

import { resolveUserName } from '../../util/resolveUserName';

import { useProfile } from '@/services/api';
import { useLocale } from '@/services/dashboard';
import { SelectedDashboardContext } from '@/services/selectedDashboardProvider';

type Page = {
  label: string;
  icon?: TablerIcon;
  path: string;
  subPages?: Page[];
};

const navigationLinks: Array<Page> = [
  { label: 'Dashboard', icon: IconDashboard, path: '/' },
  { label: 'Integrations', icon: IconBox, path: '/integrations' },
  {
    label: 'Song Requests',
    icon: IconHeadphones,
    path: '/song-requests',
    subPages: [
      { label: 'Player', icon: IconPlayerPlay, path: '/song-requests/player' },
      { label: 'Settings', icon: IconSettings, path: '/song-requests/settings' },
    ],
  },
  { label: 'Settings', icon: IconSettings, path: '/settings' },
  {
    label: 'Commands',
    icon: IconCommand,
    path: '/commands',
    subPages: [
      { label: 'Custom', icon: IconPencilPlus, path: '/commands/custom' },
      { label: 'Moderation', icon: IconSword, path: '/commands/moderation' },
      { label: 'Song Requests', icon: IconPlaylist, path: '/commands/songrequest' },
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
  opened: boolean;
  setOpened: React.Dispatch<React.SetStateAction<boolean>>;
};

const useStyles = createStyles((theme) => ({
  link: {
    borderLeft: `1px solid ${
      theme.colorScheme === 'dark' ? theme.colors.dark[4] : theme.colors.gray[3]
    }`,
  },
}));

export function SideBar(props: Props) {
  const { classes } = useStyles();
  const { locale } = useLocale();
  const viewPort = useViewportSize();
  const router = useRouter();
  const theme = useMantineTheme();
  const { t } = useTranslation('layout');

  const { data: user } = useProfile();

  const spotlight = useSpotlight();
  const dashboardContext = useContext(SelectedDashboardContext);

  const [selectedDashboard, setSelectedDashboard] = useState<Dashboard>();

  const setDefaultDashboard = useCallback(() => {
    if (!user) return;
    dashboardContext.setId(user.id);
    setSelectedDashboard({
      id: user.id,
      channelId: user.id,
      userId: user.id,
      twitchUser: user,
    });
  }, [user]);

  useEffect(() => {
    if (dashboardContext.id && user) {
      if (dashboardContext.id === user.id) {
        setDefaultDashboard();
      } else {
        const dashboard = user.dashboards.find((d) => d.channelId === dashboardContext.id);
        if (dashboard) {
          setSelectedDashboard(dashboard);
        } else {
          setDefaultDashboard();
        }
      }
    }
  }, [user, dashboardContext.id]);

  function openSpotlight() {
    if (user && dashboardContext.id) {
      spotlight.removeActions(spotlight.actions.map((a) => a.id!));
      const actions = user.dashboards
        .filter((d) => d.channelId != dashboardContext.id)
        .map((d) => ({
          title: resolveUserName(d.twitchUser.login, d.twitchUser.display_name),
          description: d.twitchUser.id,
          onTrigger: () => {
            dashboardContext.setId(d.channelId);
          },
          icon: <Avatar radius="xs" src={d.twitchUser.profile_image_url} />,
          id: d.id,
        }));

      if (dashboardContext.id != user.id) {
        actions.unshift({
          title: resolveUserName(user.login, user.display_name),
          description: user.id,
          onTrigger: () => {
            dashboardContext.setId(user.id);
          },
          icon: <Avatar radius="xs" src={user.profile_image_url} />,
          id: user.id,
        });
      }

      spotlight.registerActions(actions);
      spotlight.openSpotlight();
    }
  }

  const computeActive = (item: Page) => {
    if (item.subPages) {
      return router.asPath.startsWith(item.path);
    } else {
      return item.path === router.asPath;
    }
  };

  const createNavLink = (item: Page, isSubPage = false) => (
    <NavLink
      key={item.label}
      active={computeActive(item)}
      label={item.label}
      className={isSubPage ? classes.link : ''}
      defaultOpened={item.subPages && router.asPath.startsWith(item.path)}
      icon={item.icon ? <item.icon size={16} stroke={1.5} /> : ''}
      sx={{ width: '100%' }}
      component="a"
      href={`/dashboard/${locale}${item.path ? item.path : item.label.toLowerCase()}`}
      onClick={(e) => {
        e.preventDefault();
        if (item.subPages) return;
        props.setOpened(false);
        router.push(item.path ? item.path : item.label.toLowerCase(), undefined, { locale });
      }}
    >
      {item.subPages && item.subPages.map((p) => createNavLink(p, true))}
    </NavLink>
  );

  const [links, setLinks] = useState<JSX.Element[]>([]);
  useEffect(() => {
    setLinks(navigationLinks.map((item) => createNavLink(item)));
  }, [router]);

  return (
    <Navbar zIndex={99} hiddenBreakpoint="sm" hidden={!props.opened} width={{ sm: 250 }}>
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
                  {selectedDashboard
                    ? resolveUserName(
                        selectedDashboard.twitchUser.login,
                        selectedDashboard.twitchUser.display_name,
                      )
                    : ''}
                </Text>
              </Box>
            </Group>
          </UnstyledButton>
        </Box>
      </Navbar.Section>
    </Navbar>
  );
}
