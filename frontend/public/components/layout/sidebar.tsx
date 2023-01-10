import {
  Avatar,
  Box,
  Center,
  Group,
  Loader,
  Navbar,
  NavLink,
  ScrollArea,
  Text,
} from '@mantine/core';
import { IconCommand, IconPlaylist } from '@tabler/icons';
import { useRouter } from 'next/router';

import { useTheme } from '../../hooks/useTheme';

import { useUsersByNames } from '@/services/users';

type Props = {
  opened: boolean;
  setOpened: React.Dispatch<React.SetStateAction<boolean>>;
};

const PAGES = [
  { name: 'Commands', path: 'commands', icon: IconCommand },
  { name: 'Song Requests', path: 'song-requests', icon: IconPlaylist },
];

export const SideBar = (props: Props) => {
  const router = useRouter();
  const { theme } = useTheme();
  const { data: users, isLoading } = useUsersByNames([router.query.channelName as string]);

  const pages = PAGES.map((page) => {
    const Icon = page.icon;
    return (
      <NavLink
        key={page.name}
        component="a"
        icon={<Icon size={16} stroke={1.5} />}
        href={`/p/${router.query.channelName}/${page.path}`}
        label={page.name}
        onClick={(event) => {
          event.preventDefault();
          props.setOpened(false);
          router.push(`/${router.query.channelName}/${page.path}`);
        }}
      />
    );
  });

  return (
    <Navbar hiddenBreakpoint="sm" hidden={!props.opened} width={{ sm: 250 }}>
      <Navbar.Section grow>
        <ScrollArea.Autosize maxHeight="100%">
          <Box component={ScrollArea}>{pages}</Box>
        </ScrollArea.Autosize>
      </Navbar.Section>
      <Navbar.Section>
        <Box
          sx={{
            padding: theme.spacing.sm,
            borderTop: `1px solid ${
              theme.colorScheme === 'dark' ? theme.colors.dark[4] : theme.colors.gray[2]
            }`,

            display: 'block',
            width: '100%',
            color: theme.colorScheme === 'dark' ? theme.colors.dark[0] : theme.black,
          }}
        >
          {users && users.at(0) && (
            <Group>
              <Avatar src={users.at(0)!.profile_image_url} radius="xl" />
              <Box sx={{ flex: 1 }}>
                <Text component="span">{users.at(0)!.display_name}</Text>
              </Box>
            </Group>
          )}
          {isLoading ||
            (!users?.at(0) && (
              <Center py={'lg'}>
                <Loader variant="dots" />
              </Center>
            ))}
        </Box>
      </Navbar.Section>
    </Navbar>
  );
};
