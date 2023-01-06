import { Anchor, Button, Card, Flex, Group, List, Loader, Text } from '@mantine/core';
import {
  IconLink,
  IconPlayerPause,
  IconPlayerPlay,
  IconPlayerTrackNext,
  IconPlaylist,
  IconUser,
} from '@tabler/icons';
import YouTube from 'react-youtube';

import { resolveUserName } from '../../../util/resolveUserName';
import { PlayerSlider } from './slider';

import { usePlayer } from '@/components/song-requests/hook';
import { useCardStyles } from '@/styles/card';

const YoutubePlayer: React.FC = () => {
  const { classes: cardClasses } = useCardStyles();
  const { player, playerRef, videos, currentVideo, isPlaying, skipVideo, togglePlayState, ...options } =
    usePlayer();

  return (
    <Card withBorder radius="md" p="md" ref={playerRef}>
      <Card.Section withBorder inheritPadding py="xs">
        <Group position="apart">
          <Text weight={500}>Current Song</Text>
        </Group>
      </Card.Section>
      {options.videoId ? (
        <>
          <Card.Section className={cardClasses.card}>
            <YouTube {...options} videoId={currentVideo.videoId} />
          </Card.Section>
          <Card.Section p="md" className={cardClasses.card}>
            <Flex direction="column" gap="sm">
              <PlayerSlider isPlaying={isPlaying} player={player} />
              <Flex direction="row" gap="sm" align="center" justify="center">
                <Button
                  variant="outline"
                  disabled={videos.length === 0}
                  leftIcon={isPlaying ? <IconPlayerPause /> : <IconPlayerPlay />}
                  onClick={() => togglePlayState()}
                >
                  {isPlaying ? 'Pause' : 'Play'}
                </Button>
                <Button
                  variant="outline"
                  disabled={videos.length === 0}
                  onClick={() => skipVideo()}
                  leftIcon={<IconPlayerTrackNext />}
                >
                  Next
                </Button>
              </Flex>
              <List spacing="xs" size="sm" center>
                <List.Item icon={<IconPlaylist size={16} />}>{currentVideo.title}</List.Item>
                <List.Item icon={<IconUser size={16} />}>
                  {resolveUserName(currentVideo.orderedByName, currentVideo.orderedByDisplayName)}
                </List.Item>
                <List.Item icon={<IconLink size={16} />}>
                  <Anchor href={`https://youtu.be/${currentVideo.videoId}`} target="_blank">
                    https://youtu.be/{currentVideo.videoId}
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
