import { useResizeObserver } from '@mantine/hooks';
import { useCallback, useContext, useState } from 'react';
import { YouTubeEvent, YouTubePlayer } from 'react-youtube';
import { Options as YouTubeOptions } from 'youtube-player/dist/types';

import { PlayerContext } from '@/components/song-requests/context';

export function usePlayer() {
  const [playerRef, rect] = useResizeObserver();
  const [player, setPlayer] = useState<YouTubePlayer | null>(null);
  const { videos, skipVideo, addVideos, isPlaying, setIsPlaying, autoPlay } =
    useContext(PlayerContext);

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

  const onEnd = useCallback(() => {
    skipVideo();
  }, [videos]);

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

  const setTime = useCallback(
    (t: number) => {
      player?.seekTo(t, true);
    },
    [player],
  );

  return {
    playerRef,
    videos,
    togglePlayState,
    skipVideo,
    addVideos,
    isPlaying,
    videoId: videos[0]?.videoId ?? '',
    onReady,
    onEnd,
    onStateChange,
    getSongDuration,
    getSongCurrentTime,
    setTime,
    opts: {
      playerVars: {
        controls: 1,
        autoplay: autoPlay,
        rel: 0,
      },
      width: rect.width + 31,
      height: 300,
    } as YouTubeOptions,
  };
}
