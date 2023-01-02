import { Button, Card, Flex, Grid, Loader, Slider, Text } from '@mantine/core';
import { useViewportSize } from '@mantine/hooks';
import { IconPlayerPause, IconPlayerPlay, IconPlayerTrackNext } from '@tabler/icons';
import { useCallback, useContext, useEffect, useState } from 'react';
import YouTube, { YouTubeEvent, YouTubePlayer } from 'react-youtube';
import { Options as YouTubeOptions } from 'youtube-player/dist/types';

import { PlayerContext } from '@/components/song-requests/context';
import { formatDuration } from '@/components/song-requests/helpers';

function usePlayer() {
  const { width } = useViewportSize();
  const [player, setPlayer] = useState<YouTubePlayer | null>(null);
  const { videos, skipVideo, addVideos, isPlaying, setIsPlaying } = useContext(PlayerContext);

  const togglePlayState = useCallback(() => {
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

  const getSongDuration = useCallback(() => {
    if (!player) return 0;
    return player?.getDuration() as unknown as number;
  }, [player]);

  const getSongCurrentTime = useCallback(() => {
    if (!player) return 0;
    return player?.getCurrentTime() as unknown as number;
  }, [player]);

  const setTime = useCallback((t: number) => {
    player?.seekTo(t, true);
  }, [player]);

  return {
    videos,
    togglePlayState,
    skipVideo,
    addVideos,
    isPlaying,
    videoId: videos[0]?.videoId ?? '',
    onReady,
    onStateChange,
    getSongDuration,
    getSongCurrentTime,
    setTime,
    opts: {
      playerVars: {
        controls: 1,
        autoplay: 0,
        rel: 0,
      },
      width: width < 450 ? 330 : 450,
      height: 250,
    } as YouTubeOptions,
  };
}

const YoutubePlayer: React.FC = () => {
  const {
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

  // const progressPercentage = useMemo(() => {
  //   const result = currentTime / songDuration * 100;
  //   return Number.isNaN(result) ? 0 : result;
  // }, [currentTime, songDuration]);

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
    <Flex direction={'row'} w={options.opts.width}>
      <Card shadow="sm" p="lg" withBorder>
        {options.videoId
          ? (<>
            <Card.Section style={{ marginTop: -20, marginLeft: -30, marginRight: -30 }}>
              <YouTube {...options} onEnd={() => skipVideo()}/>
            </Card.Section>
            <Card.Section p={'md'}>
              <Flex direction={'column'} gap={'sm'}>
                <Grid align={'center'}>
                  <Grid.Col span={12}>
                    <Slider
                      value={parseInt(currentTime.toFixed(0), 10)}
                      style={{ marginLeft: 10, marginRight: 10 }}
                      label={(v) => formatDuration(v)}
                      onChange={(v) => setTime(v)}
                      max={songDuration}
                      labelAlwaysOn
                    />
                  </Grid.Col>
                  {/*<Grid.Col span={1}><Text>{formatDuration(songDuration)}</Text></Grid.Col>*/}
                </Grid>
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
                    onClick={() => togglePlayState()}
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
                <Text>Requested by {videos[0].orderedByName}</Text>
              </Flex>
            </Card.Section></>)
          : <Card.Section>
            <Flex
              style={{ width: options.opts.width, height: 300 }}
              direction={'column'}
              align={'center'}
              justify={'center'}
              gap={'sm'}
            >
              <Loader/>
              <Text>Waiting for requests</Text>
            </Flex>
          </Card.Section>
        }

      </Card>
    </Flex>
  );
};

export default YoutubePlayer;