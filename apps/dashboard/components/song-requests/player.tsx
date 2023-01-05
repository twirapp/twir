import { Badge, Button, Card, Flex, Group, Loader, Slider, Text } from '@mantine/core';
import {
  IconPlayerPause,
  IconPlayerPlay,
  IconPlayerTrackNext,
  IconPlaylist,
  IconUser,
} from '@tabler/icons';
import { useEffect, useState } from 'react';
import YouTube from 'react-youtube';

import { convertMillisToTime } from '@/components/song-requests/helpers';
import { usePlayer } from '@/components/song-requests/hook';
import { useCardStyles } from '@/styles/card';

const YoutubePlayer: React.FC = () => {
  const { classes: cardClasses } = useCardStyles();
  const {
    playerRef,
    videos,
    skipVideo,
    togglePlayState,
    isPlaying,
    getSongCurrentTime,
    getSongDuration,
    setTime,
    ...options
  } = usePlayer();
  const [currentTime, setCurrentTime] = useState(0);
  const [songDuration, setSongDuration] = useState(0);

  useEffect(() => {
    const currentTimeInterval = setInterval(() => {
      setCurrentTime(getSongCurrentTime());
    }, 500);
    const durationInterval = setInterval(() => {
      setSongDuration(getSongDuration());
    }, 500);

    return () => {
      clearInterval(currentTimeInterval);
      clearInterval(durationInterval);
    };
  }, [isPlaying]);

  return (
    <Card withBorder radius="md" p="md" ref={playerRef}>
      <Card.Section withBorder inheritPadding py="xs" className={cardClasses.card}>
        <Group position="apart">
          <Text weight={500}>Current Song</Text>
        </Group>
      </Card.Section>
      {options.videoId ? (
        <>
          <Card.Section>
            <YouTube {...options} />
          </Card.Section>
          <Card.Section p="md">
            <Flex direction="column" gap="sm">
              <Group position="apart">
                <Text size="md">00:00</Text>
                <Text size="md">{convertMillisToTime(songDuration * 1000)}</Text>
              </Group>
              <Slider
                value={parseInt(currentTime.toFixed(0), 10)}
                label={(v) => convertMillisToTime(v * 1000)}
                onChange={(v) => setTime(v)}
                max={songDuration}
              />
              <Flex direction={'row'} gap={'sm'} align={'center'} justify={'center'}>
                <Button
                  variant={'outline'}
                  disabled={videos.length === 0}
                  leftIcon={isPlaying ? <IconPlayerPause /> : <IconPlayerPlay />}
                  onClick={() => togglePlayState()}
                >
                  {isPlaying ? 'Pause' : 'Play'}
                </Button>
                <Button
                  variant={'outline'}
                  disabled={videos.length === 0}
                  onClick={() => skipVideo()}
                  leftIcon={<IconPlayerTrackNext />}
                >
                  Next
                </Button>
              </Flex>
              <Flex direction={'column'} gap={'sm'} justify={'flex-start'}>
                <Text>
                  <IconPlaylist size={16} /> {videos[0].title}{' '}
                  <Badge>{convertMillisToTime(videos[0].duration)}</Badge>
                </Text>
                <Text>
                  <IconUser size={16} /> Requested by <Badge>{videos[0].orderedByName}</Badge>
                </Text>
              </Flex>
            </Flex>
          </Card.Section>
        </>
      ) : (
        <Card.Section>
          <Flex
            style={{ width: options.opts.width, height: 300 }}
            direction={'column'}
            align={'center'}
            justify={'center'}
            gap={'sm'}
          >
            <Loader variant="dots" />
          </Flex>
        </Card.Section>
      )}
    </Card>
  );
};

export default YoutubePlayer;
