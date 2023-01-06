import {
  ActionIcon,
  Anchor,
  Button,
  Card,
  Flex,
  Group,
  List,
  Loader,
  Text,
  Tooltip,
} from '@mantine/core';
import {
  IconPlayerPause,
  IconPlayerPlay,
  IconPlayerTrackNext,
  IconPlaylist,
  IconEyeOff,
  IconUser,
  IconEye,
  IconLink,
} from '@tabler/icons';
import { useState } from 'react';
import YouTube from 'react-youtube';

import { resolveUserName } from '../../../util/resolveUserName';
import { PlayerDurationSlider } from './duration-slider';

import { usePlayer } from '@/components/song-requests/hook';
import { useCardStyles } from '@/styles/card';

const YoutubePlayer: React.FC = () => {
  const { classes: cardClasses } = useCardStyles();
  const {
    player,
    // playerRef,
    videos,
    currentVideo,
    isPlaying,
    skipVideo,
    togglePlayState,
    ...options
  } = usePlayer();

  const [isVisiblePlayer, toggleVisiblePlayer] = useState(true);

  return (
    <Card withBorder radius="md" p="md" style={{ paddingBottom: 'initial' }}>
      <Card.Section withBorder inheritPadding py="xs">
        <Group position="apart">
          <Text weight={500}>Current Song</Text>

          <Tooltip
            withinPortal
            position="top"
            label={isVisiblePlayer ? 'Hide video player' : 'Show video player'}
          >
            <ActionIcon onClick={() => toggleVisiblePlayer((v) => !v)}>
              {isVisiblePlayer ? <IconEyeOff size={14} /> : <IconEye size={14} />}
            </ActionIcon>
          </Tooltip>
        </Group>
      </Card.Section>
      {currentVideo?.videoId ? (
        <>
          <Card.Section
            style={{ display: isVisiblePlayer ? 'block' : 'none' }}
            className={cardClasses.card}
          >
            <YouTube {...options} videoId={currentVideo.videoId} />
          </Card.Section>
          <Card.Section p="md" className={cardClasses.card}>
            <Flex direction="column" gap="sm">
              <PlayerDurationSlider isPlaying={isPlaying} player={player} />
              <Button.Group>
                <Button
                  fullWidth={true}
                  variant="default"
                  disabled={videos.length === 0}
                  onClick={() => togglePlayState()}
                  leftIcon={
                    isPlaying ? <IconPlayerPause size={16} /> : <IconPlayerPlay size={16} />
                  }
                >
                  {isPlaying ? 'Pause' : 'Play'}
                </Button>
                <Button
                  fullWidth={true}
                  variant="default"
                  disabled={videos.length === 0}
                  onClick={() => skipVideo()}
                  leftIcon={<IconPlayerTrackNext size={16} />}
                >
                  Next
                </Button>
              </Button.Group>
              <List spacing="xs" size="sm" center>
                <List.Item icon={<IconPlaylist size={16} />}>{currentVideo.title}</List.Item>
                <List.Item icon={<IconUser size={16} />}>
                  {resolveUserName(currentVideo.orderedByName, currentVideo.orderedByDisplayName)}
                </List.Item>
                <List.Item icon={<IconLink size={16} />}>
                  <Anchor href={`https://youtu.be/${currentVideo.videoId}`} target="_blank">
                    youtu.be/{currentVideo.videoId}
                  </Anchor>
                </List.Item>
              </List>
            </Flex>
          </Card.Section>
        </>
      ) : (
        <Card.Section className={cardClasses.card}>
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
