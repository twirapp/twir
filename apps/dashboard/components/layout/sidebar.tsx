import { Box, Navbar, NavLink, ScrollArea } from '@mantine/core';
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
import { useRouter } from 'next/router';
import { MutableRefObject, useState } from 'react';

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

export function SideBar({
  opened,
  reference,
}: {
  opened: boolean;
  reference: MutableRefObject<any>;
}) {
  const router = useRouter();

  const items = navigationLinks.map((item, index) => (
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
    <Navbar ref={reference} hiddenBreakpoint="sm" hidden={!opened} width={{ sm: 150, lg: 150 }}>
      <Box component={ScrollArea}>{items}</Box>
    </Navbar>
  );
}
