import { Grid } from '@mantine/core';
import { useListState } from '@mantine/hooks';
import type { RequestedSong } from '@tsuwari/typeorm/entities/RequestedSong';
import { getCookie } from 'cookies-next';
import { GetServerSideProps, NextPage } from 'next';
import { serverSideTranslations } from 'next-i18next/serverSideTranslations';
import dynamic from 'next/dynamic';
import { useCallback, useEffect, useRef, useState } from 'react';
import { DraggableLocation } from 'react-beautiful-dnd';
import { io, Socket } from 'socket.io-client';

import { PlayerContext } from '@/components/song-requests/context';
import { moveItem } from '@/components/song-requests/helpers';
import { QueueList } from '@/components/song-requests/queue/queue-list';
import { useProfile } from '@/services/api';

export const getServerSideProps: GetServerSideProps = async ({ locale }) => ({
  props: {
    ...(await serverSideTranslations(locale!, ['layout'])),
  },
});

const YoutubePlayer = dynamic(
  () => import('../../components/song-requests/player/youtube-player'),
  {
    ssr: false,
  },
);

const Player: NextPage = () => {
  const profile = useProfile();
  const [videos, videosHandlers] = useListState<RequestedSong>([]);
  const socketRef = useRef<Socket | null>(null);
  const [autoPlay, setAutoPlay] = useState(0);
  const [isPlaying, setIsPlaying] = useState(false);

  const skipVideo = useCallback(
    (index = 0, callWs = true) => {
      if (videos[index] && callWs) {
        callWsSkip(videos[index]!);
      }

      const length = videos.length;
      if (index === 0) {
        setIsPlaying(false);
        videosHandlers.setState(videos.slice(1));
      } else if (index === length - 1) {
        videosHandlers.setState(videos.slice(0, length - 1));
      } else {
        videosHandlers.setState([...videos.slice(0, index), ...videos.slice(index + 1)]);
      }
    },
    [videos],
  );

  const clearQueue = useCallback(() => {
    callWsSkip(videos);
  }, [videos]);

  const addVideos = useCallback(
    (v: RequestedSong[]) => {
      videosHandlers.setState((existedVideos) => [...existedVideos, ...v]);
    },
    [videos],
  );

  const reorderVideos = useCallback(
    (destination: DraggableLocation, source: DraggableLocation) => {
      const from = source.index;
      const to = destination?.index || 0;
      videosHandlers.reorder({ from, to });

      const newVideos = moveItem(videos, from, to)
        .map((v, i) => ({ ...v, queuePosition: i + 1 }))
        .slice(1);

      socketRef.current?.emit('newOrder', newVideos);
    },
    [videos, socketRef.current],
  );

  useEffect(() => {
    if (!profile.data) return;

    if (!socketRef.current) {
      socketRef.current = io(
        `${`${window.location.protocol == 'https:' ? 'wss' : 'ws'}://${
          window.location.host
        }`}/youtube`,
        {
          transports: ['websocket'],
          autoConnect: false,
          auth: (cb) => {
            cb({ apiKey: profile.data?.apiKey, channelId: getCookie('selectedDashboard') });
          },
        },
      );
    }

    socketRef.current.connect();

    socketRef.current.emit('currentQueue', (data: RequestedSong[]) => {
      videosHandlers.setState([]);
      addVideos(data);
    });

    socketRef.current.on('newTrack', (track: RequestedSong) => {
      if (autoPlay === 0) setAutoPlay(1);
      addVideos([track]);
    });

    return () => {
      socketRef.current?.off('newTrack');
      socketRef.current?.disconnect();
    };
  }, [profile.data]);

  // it's in another useEffect because we need videos as dependency for correctly find index for skip
  useEffect(() => {
    if (!socketRef.current) return;

    socketRef.current.on('removeTrack', (track: RequestedSong) => {
      const index = videos.findIndex((v) => v.id === track.id);
      if (index > 0) {
        skipVideo(index, false);
      }
    });

    return () => {
      socketRef.current?.off('removeTrack');
    };
  }, [videos, socketRef.current]);

  function callWsSkip(videos: RequestedSong | RequestedSong[]) {
    const ids = (Array.isArray(videos) ? videos : [videos]).map((v) => v.id);
    socketRef.current?.emit('skip', ids);
  }

  useEffect(() => {
    const video = videos[0]!;
    if (isPlaying) {
      socketRef.current?.emit('play', { id: video.id, duration: video.duration });
    } else {
      socketRef.current?.emit('pause');
    }
  }, [isPlaying]);

  return (
    <Grid>
      <PlayerContext.Provider
        value={{
          videos,
          videosHandlers,
          skipVideo,
          clearQueue,
          addVideos,
          isPlaying,
          setIsPlaying,
          reorderVideos,
          autoPlay,
          setAutoPlay,
        }}
      >
        <Grid.Col md={4} lg={4}>
          <YoutubePlayer />
        </Grid.Col>
        <Grid.Col md={8} lg={8}>
          <QueueList />
        </Grid.Col>
      </PlayerContext.Provider>
    </Grid>
  );
};

export default Player;
