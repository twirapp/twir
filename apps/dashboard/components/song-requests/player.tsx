import { Button, Flex, Text } from '@mantine/core';
import { IconPlayerPause, IconPlayerPlay, IconPlayerTrackNext } from '@tabler/icons';
import { useCallback, useContext, useState } from 'react';
import YouTube, { YouTubeEvent, YouTubePlayer } from 'react-youtube';
import { Options as YouTubeOptions } from 'youtube-player/dist/types';

import { PlayerContext } from '@/components/song-requests/context';

function usePlayer() {
  const [player, setPlayer] = useState<YouTubePlayer | null>(null);
  const { videos, skipVideo, addVideos, isPlaying, setIsPlaying } = useContext(PlayerContext);

  const toggleVideo = useCallback(() => {
    if (isPlaying) {
      player?.pauseVideo();
    } else {
      player?.playVideo();
    }
  }, [player, isPlaying]);

  const onReady = useCallback((event: YouTubeEvent) => {
    setPlayer(event.target);
  }, []);

  const onStateChange = useCallback(
    (event: YouTubeEvent<number>) => {
      switch (event.data) {
        case 1:
          setIsPlaying(true);
          break;
        case 2:
          setIsPlaying(false);
          break;
        case -1:
          player?.playVideo();
          break;
      }
    },
    [player],
  );

  return {
    videos,
    toggleVideo,
    skipVideo,
    addVideos,
    isPlaying,
    videoId: videos[0]?.videoId ?? '',
    onReady,
    onStateChange,
    opts: {
      playerVars: {
        controls: 1,
        autoplay: 0,
        rel: 0,
      },
      width: 450,
      height: 250,
    } as YouTubeOptions,
  };
}

const YoutubePlayer: React.FC = () => {
  const { videos, skipVideo, isPlaying, ...options } = usePlayer();

  return (
    <Flex direction={'column'} gap={'md'} w={options.opts.width}>
      {options.videoId
        ? <YouTube {...options} onEnd={() => skipVideo()}/>
        : <Text size={'xl'}>Waiting for songs...</Text>
      }
      <Flex
        direction={'row'}
        gap={'sm'}
        align={'center'}
        justify={'center'}
      >
        <Button
          variant={'outline'}
          disabled={videos.length === 0}
          leftIcon={isPlaying ? <IconPlayerPause/> : <IconPlayerPlay/>}
        >
          {isPlaying ? 'Pause' : 'Play'}
        </Button>
        <Button
          variant={'outline'}
          disabled={videos.length === 0}
          onClick={() => skipVideo()}
          leftIcon={<IconPlayerTrackNext/>}
        >
          Next
        </Button>

      </Flex>
    </Flex>
  );
};

export default YoutubePlayer;