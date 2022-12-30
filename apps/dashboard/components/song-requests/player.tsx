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
    } as YouTubeOptions,
  };
}

const YoutubePlayer: React.FC = () => {
  const { videos, toggleVideo, skipVideo, addVideos, isPlaying, ...options } = usePlayer();

  return (
    <div>
      {options.videoId ? <YouTube {...options} onEnd={() => skipVideo()}
      /> : <h1>Queue is empty</h1>}
      <button type="button" onClick={() => toggleVideo()}>
        {isPlaying ? 'Pause' : 'Play'}
      </button>
      <button type="button" onClick={() => skipVideo()}>
        Next
      </button>
    </div>
  );
};

export default YoutubePlayer;