import { Button } from '@mantine/core';
import type { RequestedSong } from '@tsuwari/typeorm/entities/RequestedSong';
import { getCookie } from 'cookies-next';
import { useCallback, useContext, useEffect, useRef, useState } from 'react';
import YouTube, { YouTubeEvent } from 'react-youtube';
import { io, Socket } from 'socket.io-client';
import type { Options as YouTubePlayerOptions, YouTubePlayer } from 'youtube-player/dist/types';

import { PlayerContext } from '@/components/song-requests/context';
import { useProfile } from '@/services/api';

const baseWsUrl = `${window.location.protocol == 'https:' ? 'wss' : 'ws'}://${window.location.host}`;

type UsePlayerProps = {
  onSkip: (track: RequestedSong) => void
  onPlay: (id: string, timeToEnd: number) => void
  onPause: () => void
}

export function usePlayer(props: UsePlayerProps) {
  const { videos, setVideos } = useContext(PlayerContext);

  const [isPlaying, setIsPlaying] = useState(false);
  const [player, setPlayer] = useState<YouTubePlayer | null>(null);

  const emitSkip = useCallback(() => {
    const current = videos[0];
    if (current) {
      props.onSkip(current);
    }
  }, [videos]);

  const next = useCallback(() => {
    emitSkip();

    setVideos(videos => videos.slice(1));

    if (videos.length) {
      player?.seekTo(0, false);
      player?.playVideo();
    }
  }, [videos]);

  const toggle = useCallback(() => {
    if (isPlaying) {
      player?.pauseVideo();
    } else {
      player?.playVideo();
    }
  }, [isPlaying, player]);

  const setPlayerVideos = (data: RequestedSong[]) => {
    setVideos(videos => data);
  };

  const onReady = useCallback(
    (event: YouTubeEvent<any>) => {
      setPlayer(event.target);
    },
    [player, videos],
  );

  // -1 (воспроизведение видео не начато)
  // 0 (воспроизведение видео завершено)
  // 1 (воспроизведение)
  // 2 (пауза)
  // 3 (буферизация)
  // 5 (видео подают реплики).
  const onStateChange = useCallback(
    (event: YouTubeEvent<any>) => {
      if (event.data === 0) {
        setIsPlaying(false);
        next();
      } else if (event.data === 1) {
        setIsPlaying(true);
        const currentTime = player?.getCurrentTime() as unknown as number;
        const video = videos[0]!;
        const timeToEnd = video.duration - (Math.round(currentTime) * 1000);
        props.onPlay(videos[0]!.id, timeToEnd);
      } else if (event.data === 2) {
        setIsPlaying(false);
        props.onPause();
      }
    },
    [player],
  );

  return {
    toggle,
    next,
    isPlaying,
    videoId: videos[0]?.videoId ?? '',
    onReady,
    onStateChange,
    setPlayerVideos,
    opts: {
      playerVars: {
        controls: 1,
        autoplay: 0,
        modestbranding: 1,
        showinfo: 0,
        rel: 0,
        ecver: 2,
        loop: 0,
      },
    } as YouTubePlayerOptions,
  };
}

const YoutubePlayer: React.FC = () => {
  const { videos } = useContext(PlayerContext);

  const { toggle, next, isPlaying, setPlayerVideos, ...options } = usePlayer({
    onSkip,
    onPlay,
    onPause,
  });
  const profile = useProfile();
  const socketRef = useRef<Socket | null>(null);

  function onSkip(track: RequestedSong) {
    console.log('skiping');
    // socketRef.current?.emit('skip', track.id);
    return;
  }

  function onPlay(id: string, timeToEnd: number) {
    socketRef.current?.emit('play', { id, timeToEnd });
  }

  function onPause() {
    socketRef.current?.emit('pause');
  }

  useEffect(() => {
    if (!socketRef.current) {
      socketRef.current = io(`${baseWsUrl}/youtube`, {
        transports: ['websocket'],
        autoConnect: false,
        auth: (cb) => {
          cb({ apiKey: profile.data?.apiKey, channelId: getCookie('selectedDashboard') });
        },
      });
    }

    socketRef.current.connect();

    socketRef.current.emit('currentQueue', (data: RequestedSong[]) => {
      setPlayerVideos(data);
    });

    // TODO: unsubscribe from events
    return () => {
      socketRef.current?.disconnect();
    };
  }, [profile.data]);

  return <div>
    <Button onClick={toggle}>{JSON.stringify(isPlaying)}</Button>
    <Button onClick={next} disabled={videos.length === 0}>next</Button>
    <YouTube {...options} />
  </div>;
};

export default YoutubePlayer;