import {
  ActionIcon,
  Box,
  Button,
  Card,
  Center,
  Divider,
  Flex,
  Grid,
  Group,
  Menu,
  ScrollArea,
  Text,
  Transition,
} from '@mantine/core';
import { useListState } from '@mantine/hooks';
import {
  IconAdjustmentsHorizontal,
  IconCircleMinus,
  IconPlayerPlay,
  IconPlayerSkipForward,
  IconPlaylist,
  IconVideoOff,
} from '@tabler/icons';
import Plyr, { APITypes, PlyrOptions } from 'plyr-react';
import React, { useEffect, useRef, useState } from 'react';
import 'plyr-react/plyr.css';
import { useRouter } from 'next/router';

const plyrOptions: PlyrOptions = {
  controls: [
    'progress',
    'current-time',
    'mute',
    'volume',
    'captions',
    'settings',
    'pip',
    'airplay',
    'fullscreen',
  ],
  ratio: '16:9',
  hideControls: true,
  keyboard: { focused: false, global: false },
  invertTime: false,
  debug: false,
};

type Track = Plyr.Source & {
  title: string,
  orderedBy: string
}

export const YoutubePlayer: React.FC = () => {
  const router = useRouter();
  const plyrRef = useRef<APITypes>(null) as React.MutableRefObject<APITypes>;

  const [currentTrack, setCurrentTrack] = useState<Track>();
  const [songs, songsHandlers] = useListState<Track>([
    {
      src: 'WLcHVVS90zQ',
      provider: 'youtube',
      title: 'Test',
      orderedBy: 'Satont',
    },
    {
      src: 'FCtasDPQ9e8',
      provider: 'youtube',
      title: 'Test 2',
      orderedBy: 'mellkam',
    }, {
      src: 'FCtasDPQ9e8',
      provider: 'youtube',
      title: 'Test 2',
      orderedBy: 'mellkam',
    }, {
      src: 'FCtasDPQ9e8',
      provider: 'youtube',
      title: 'Test 2',
      orderedBy: 'mellkam',
    }, {
      src: 'FCtasDPQ9e8',
      provider: 'youtube',
      title: 'Test 2',
      orderedBy: 'mellkam',
    }, {
      src: 'FCtasDPQ9e8',
      provider: 'youtube',
      title: 'Test 2',
      orderedBy: 'mellkam',
    }, {
      src: 'FCtasDPQ9e8',
      provider: 'youtube',
      title: 'Test 2',
      orderedBy: 'mellkam',
    },
  ]);

  useEffect(() => {
    if (songs.at(0)) {
      setCurrentTrack(songs.at(0));
    } else {
      setCurrentTrack(undefined);
    }
  }, [songs]);

  const nextVideo = () => {
    songsHandlers.shift();
  };

  return <Grid grow>
    <Grid.Col span={4}>
      <Card>
        <Card.Section p={'xs'}>
          <Flex gap="xs" direction="row" justify="space-between">
            <Text size="md">YouTube</Text>

            <Group>
              <Menu shadow="md" width={400} styles={{ dropdown: { backgroundColor: '#2C2C2C' } }}>
                <Menu.Target>
                  <ActionIcon><IconPlaylist /></ActionIcon>
                </Menu.Target>

                <Menu.Dropdown h={200}>
                  <ScrollArea h={190} type={'auto'}>
                    {songs.map(s => <Menu.Item key={s.src} rightSection={
                      <ActionIcon><IconCircleMinus /></ActionIcon>
                    }>
                      <Flex direction={'column'}>
                        <Text size={'lg'}>{s.title}</Text>
                        <Text size={'xs'}>Ordered by: {s.orderedBy}</Text>
                      </Flex>
                    </Menu.Item>)}
                  </ScrollArea>
                </Menu.Dropdown>
              </Menu>
              <ActionIcon onClick={() => router.push('/settings/youtube')}><IconAdjustmentsHorizontal /></ActionIcon>
            </Group>
          </Flex>
        </Card.Section>
        <Divider />
        <Card.Section>
          <Box sx={{ height: 287 }} hidden={!!songs.length}>
            <Center style={{ height: 287 }}>
              <Flex direction={'column'} align={'center'}>
                <IconVideoOff size={130} />
                <Text>there is no video in queue</Text>
              </Flex>
            </Center>
          </Box>
          <div hidden={!songs.length}>
            <Plyr
              ref={plyrRef as any}
              source={{
                type: 'video',
                sources: currentTrack ? [currentTrack] : [],
              }}
              options={plyrOptions}
            />
          </div>
        </Card.Section>
        <Transition mounted={!!currentTrack} transition="slide-down" duration={200} timingFunction="ease">
          {(styles) => <Card.Section p={'xs'} style={styles}>
            <Flex direction={'row'} justify={'space-between'}>
              <Flex direction={'column'}>
                <Text size={'lg'}>{currentTrack?.title}</Text>
                <Text size={'xs'} color={'lime'}>Ordered by: {currentTrack?.orderedBy}</Text>
              </Flex>
              <Group>
                <ActionIcon><IconPlayerPlay /></ActionIcon>
                <ActionIcon><IconPlayerSkipForward onClick={nextVideo} /></ActionIcon>
              </Group>
            </Flex>
          </Card.Section>}
        </Transition>
      </Card>
    </Grid.Col>
  </Grid>;
};
