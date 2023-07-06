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
  Button,
  Divider,
  Loader,
} from '@mantine/core';
import { useViewportSize } from '@mantine/hooks';
import { useSpotlight } from '@mantine/spotlight';
import {
  IconActivity,
  IconBox,
  IconCalendarEvent,
  IconClipboardCopy,
  IconClockHour7,
  IconCommand,
  IconDashboard,
  IconDeviceDesktop,
  IconDeviceDesktopAnalytics,
  IconExternalLink,
  IconHeadphones,
  IconKey,
  IconPencilPlus,
  IconPlayerPlay,
  IconPlaylist,
  IconSettings,
  IconShieldHalfFilled,
  IconSpeakerphone,
  IconSword,
  IconUsers,
  TablerIcon,
} from '@tabler/icons';
import { Dashboard } from '@twir/grpc/generated/api/api/auth';
import { useTranslation } from 'next-i18next';
import { useRouter } from 'next/router';
import { useCallback, useContext, useEffect, useState } from 'react';

import { useDashboards, useProfile } from '@/services/api';
import { useLocale } from '@/services/dashboard';
import { SelectedDashboardContext } from '@/services/selectedDashboardProvider';


type Page = {
  label: string;
  icon?: TablerIcon;
  path: string;
  subPages?: Page[];
};

const navigationLinks: Array<Page | null> = [
  { label: 'Dashboard', icon: IconDashboard, path: '/' },
  { label: 'Integrations', icon: IconBox, path: '/integrations' },
  { label: 'Events', icon: IconCalendarEvent, path: '/events' },
  { label: 'Overlays', icon: IconDeviceDesktop, path: '/overlays' },
  {
    label: 'Song Requests',
    icon: IconHeadphones,
    path: '/song-requests',
    subPages: [
      { label: 'Player', icon: IconPlayerPlay, path: '/song-requests/player' },
      { label: 'Settings', icon: IconSettings, path: '/song-requests/settings' },
    ],
  },
  {
    label: 'Commands',
    icon: IconCommand,
    path: '/commands',
    subPages: [
      { label: 'Custom', icon: IconPencilPlus, path: '/commands/custom' },
      { label: 'Stats', icon: IconDeviceDesktopAnalytics, path: '/commands/stats' },
      // { label: 'Moderation', icon: IconSword, path: '/commands/moderation' },
      { label: 'Songs', icon: IconPlaylist, path: '/commands/songs' },
      { label: 'Manage', icon: IconClipboardCopy, path: '/commands/manage' },
      // { label: 'Dota', path: '/commands/dota' },
    ],
  },
  {
    label: 'Community',
    icon: IconUsers,
    path: '/community',
    subPages: [
      { label: 'Users', icon: IconUsers, path: '/community/users' },
      { label: 'Roles', icon: IconShieldHalfFilled, path: '/community/roles' },
    ],
  },
  { label: 'Timers', icon: IconClockHour7, path: '/timers' },
  { label: 'Moderation', icon: IconSword, path: '/moderation' },
  { label: 'Keywords', icon: IconKey, path: '/keywords' },
  { label: 'Variables', icon: IconActivity, path: '/variables' },
  { label: 'Greetings', icon: IconSpeakerphone, path: '/greetings' },
  null,
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
  const { data: dashboards, isLoading } = useDashboards();

  const [selectedDashboard, setSelectedDashboard] = useState<Dashboard>();
  const spotlight = useSpotlight();
  const dashboardContext = useContext(SelectedDashboardContext);

  useEffect(() => {
    if (!user || !dashboards || !dashboardContext.id) return;
    const dashboard = dashboards.dashboards.find((d) => d.id === dashboardContext.id);
    setSelectedDashboard(dashboard);
  }, [user, dashboards, dashboardContext.id]);

  const openSpotlight = useCallback(() => {
    if (!dashboards) return;

    if (user && dashboardContext.id) {
      spotlight.removeActions(spotlight.actions.map((a) => a.id!));

      const actions = dashboards.dashboards.map((d) => ({
        title: d.id,
        description: d.id,
        onTrigger: () => {
          setSelectedDashboard(dashboards.dashboards.find((dash) => dash.id === d.id));
          dashboardContext.setId(d.id);
        },
        icon: <Avatar radius="xs" src={d.id} />,
        id: d.id,
      }));

      spotlight.registerActions(actions);
      spotlight.openSpotlight();
    }
  }, [dashboards]);

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
    setLinks(
      navigationLinks.map((item, i) => {
        return item === null ? <Divider key={i} /> : createNavLink(item);
      }),
    );
  }, [router]);

  return (
    <Navbar zIndex={99} hiddenBreakpoint="sm" hidden={!props.opened} width={{ sm: 250 }}>
      <Navbar.Section grow>
        <ScrollArea.Autosize
          maxHeight={viewPort.height - 170}
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
          {!isLoading && (
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
                <Avatar src={''} radius="xl" />
                <Box sx={{ flex: 1 }}>
                  <Text size="xs" weight={500}>
                    {t('sidebar.manage')}
                  </Text>
                  <Text color="dimmed" size="xs">
                    {selectedDashboard
                      ? 'here dashboard twitch user name'
                      : ''}
                  </Text>
                </Box>
              </Group>
            </UnstyledButton>
          )}
          {isLoading && <Loader color="violet" variant="dots" size={50} />}
          <Button
            size={'xs'}
            compact
            color="grape"
            style={{ marginTop: 5 }}
            variant={'light'}
            component="a"
            href={
              // 'window' in globalThis && selectedDashboard?.name
              //   ? `${window.location.origin}/p/${selectedDashboard?.name}/commands`
              //   : ''
							''
            }
            target={'_blank'}
            leftIcon={<IconExternalLink size={14} />}
            w={'100%'}
          >
            Public page
          </Button>
        </Box>
      </Navbar.Section>
    </Navbar>
  );
}
