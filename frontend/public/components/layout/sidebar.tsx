import {
  ActionIcon,
  Avatar,
  Box,
  Center,
  Flex,
  Group,
  Loader,
  Navbar,
  NavLink,
  ScrollArea,
  Text,
  useMantineColorScheme,
  useMantineTheme,
} from '@mantine/core';
import { useHotkeys, useViewportSize } from '@mantine/hooks';
import { IconCommand, IconMoonStars, IconPlaylist, IconSun } from '@tabler/icons';
import Image from 'next/image';
import { useRouter } from 'next/router';

import { useUsersByNames } from '@/services/users';

type Props = {
  opened: boolean;
  setOpened: React.Dispatch<React.SetStateAction<boolean>>;
};

export const useTheme = () => {
  const theme = useMantineTheme();
  const { colorScheme, toggleColorScheme } = useMantineColorScheme();

  useHotkeys([['mod+J', () => toggleColorScheme()]]);

  return {
    theme,
    colorScheme,
    toggleColorScheme,
  };
};

export const SideBar = (props: Props) => {
  const viewPort = useViewportSize();
  const router = useRouter();
  const { theme, colorScheme, toggleColorScheme } = useTheme();

  const { data: users, isLoading } = useUsersByNames([router.query.channelName as string]);

  const links = [
    { name: 'Commands', path: 'commands', icon: IconCommand },
    { name: 'Song requests', path: 'song-requests', icon: IconPlaylist },
  ].map((l) => {
    return <NavLink
      key={l.name}
      icon={<l.icon size={16} stroke={1.5}/>}
      component="a"
      href={`/p/${router.query.channelName}/${l.path}`}
      label={l.name}
      onClick={(e) => {
        e.preventDefault();
        props.setOpened(false);
        router.push(`/${router.query.channelName}/${l.path}`);
      }}
    />;
  });

  return <Navbar zIndex={99} hiddenBreakpoint="sm" hidden={!props.opened} width={{ sm: 250 }}>
    <Navbar.Section>
        <Box
            sx={{
              padding: theme.spacing.sm,
              borderBottom: `1px solid ${
                theme.colorScheme === 'dark' ? theme.colors.dark[4] : theme.colors.gray[2]
              }`,

              display: 'block',
              width: '100%',
              color: theme.colorScheme === 'dark' ? theme.colors.dark[0] : theme.black,
            }}
        >
          {users && users.at(0) &&
            <Group>
                <Avatar src={users.at(0)!.profile_image_url} radius="xl"/>
                <Box sx={{ flex: 1 }}>
                    <Text component="span">
                      {users.at(0)!.display_name}
                    </Text>
                </Box>
            </Group>
          }
          {isLoading || !users?.at(0) && <Center py={'lg'}>
              <Loader variant="dots" />
          </Center>}
        </Box>
    </Navbar.Section>
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

          display: 'block',
          width: '100%',
          color: theme.colorScheme === 'dark' ? theme.colors.dark[0] : theme.black,
        }}
      >
        <Flex align={'center'} justify={'space-between'}>
          <Group>
            <Image src="/p/TsuwariInCircle.svg" width={30} height={30} alt="Tsuwari Logo"/>
            <Text
              component="span"
            >
              Tsuwari
            </Text>
          </Group>
          <ActionIcon
            size="lg"
            variant="default"
            color={colorScheme === 'dark' ? 'yellow' : 'blue'}
            onClick={() => toggleColorScheme()}
            title="Toggle color scheme"
          >
            {colorScheme === 'dark' ? <IconSun size={18}/> : <IconMoonStars size={18}/>}
          </ActionIcon>
        </Flex>
      </Box>
    </Navbar.Section>
  </Navbar>;
};