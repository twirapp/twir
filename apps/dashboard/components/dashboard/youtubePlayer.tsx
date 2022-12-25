import { ActionIcon, Button, Card, Divider, Flex, Grid, Group, Text } from '@mantine/core';
import { IconPlayerPlay, IconPlayerSkipForward } from '@tabler/icons';
import Plyr, { APITypes, PlyrSource, PlyrOptions } from 'plyr-react';
import React, { useRef, useState } from 'react';
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

export const YoutubePlayer: React.FC = () => {
  const plyrRef = useRef<APITypes>(null) as React.MutableRefObject<APITypes>;

  const [source, setSource] = useState<PlyrSource>( {
    type: 'video',
    sources: [
      {
        src: 'WLcHVVS90zQ',
        provider: 'youtube',
      },
    ],
  });

  const setVideo = () => {
    plyrRef.current.plyr.source = {
      type: 'video',
      sources: [
        {
          src: 'FCtasDPQ9e8',
          provider: 'youtube',
        },
      ],
    };
  };

  return <Grid grow>
    <Grid.Col span={4}>
      <Card>
        <Card.Section p={'xs'}>
          <Flex gap="xs" direction="row" justify="space-between">
            <Text size="md">YouTube</Text>
          </Flex>
        </Card.Section>
        <Divider />
        <Card.Section>
          <Plyr
            ref={plyrRef as any}
            source={source}
            options={plyrOptions}
          />
          {/*<video ref={ref} className="plyr-react plyr" {...rest} />*/}
        </Card.Section>
        <Card.Section p={'xs'}>
          <Flex direction={'row'} justify={'space-between'}>
            <Flex direction={'column'}>
              <Text size={'lg'}>qweqweqweqweqweqweqweqweqweqweqwe</Text>
              <Text size={'xs'} color={'lime'}>Ordered by: mellkam</Text>
            </Flex>
            <Group>
              <ActionIcon><IconPlayerPlay /></ActionIcon>
              <ActionIcon><IconPlayerSkipForward /></ActionIcon>
            </Group>
          </Flex>
        </Card.Section>
      </Card>
    </Grid.Col>
  </Grid>;
};
