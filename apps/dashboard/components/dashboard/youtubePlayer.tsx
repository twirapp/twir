import { ActionIcon, Box, Button, Card, Center, Divider, Flex, Grid, Group, Text, Transition } from '@mantine/core';
import { useListState } from '@mantine/hooks';
import { IconPlayerPlay, IconPlayerSkipForward, IconVideoOff } from '@tabler/icons';
import Plyr, { APITypes, PlyrOptions } from 'plyr-react';
import React, { useEffect, useRef, useState } from 'react';
import 'plyr-react/plyr.css';

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
    },
  ]);

  useEffect(() => {
    if (songs.at(0)) {
      setCurrentTrack(songs.at(0));
    } else {
      setCurrentTrack(undefined);
    }
  }, [songs]);

  const setVideo = () => {
    songsHandlers.shift();
  };

  return <Grid grow>
    <Grid.Col span={4}>
      <Card>
        <Card.Section p={'xs'}>
          <Flex gap="xs" direction="row" justify="space-between">
            <Text size="md">YouTube</Text>
            {/*<Button onClick={setVideo} size={'xs'}>test change source</Button>*/}
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
                <ActionIcon><IconPlayerSkipForward /></ActionIcon>
              </Group>
            </Flex>
          </Card.Section>}
        </Transition>
      </Card>
    </Grid.Col>
  </Grid>;
};
