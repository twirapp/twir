import { Box, Navbar, NavLink, ScrollArea, useMantineTheme } from '@mantine/core';
import { useViewportSize } from '@mantine/hooks';
import { IconCommand, IconPlaylist } from '@tabler/icons';
import { useRouter } from 'next/router';

type Props = {
  opened: boolean;
  setOpened: React.Dispatch<React.SetStateAction<boolean>>;
};

export const SideBar = (props: Props) => {
  const viewPort = useViewportSize();
  const router = useRouter();
  const theme = useMantineTheme();

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
        router.push(l.path, { query: { channelName: router.query.channelName } });
      }}
    />;
  });

  return <Navbar zIndex={99} hiddenBreakpoint="sm" hidden={!props.opened} width={{ sm: 250 }}>
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
  </Navbar>;
};