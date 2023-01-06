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

  return {
    player,
    playerRef,
    videos,
    currentVideo: videos[0],
    isPlaying,
    videoId: videos[0]?.videoId ?? '',
    togglePlayState,
    skipVideo,
    addVideos,
    // ...options
    onReady,
    onEnd,
    onStateChange,
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
