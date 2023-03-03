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
import { AlertPlayerDisabled, AlertQueueEmpty } from '@/components/song-requests/player/alerts';
import { QueueList } from '@/components/song-requests/queue/queue-list';
import { useProfile } from '@/services/api';
import { useYoutubeModule } from '@/services/api/modules';

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
  const youtube = useYoutubeModule();
  const { data: youtubeSettings } = youtube.useSettings();
  const profile = useProfile();
  const [videos, videosHandlers] = useListState<RequestedSong>([]);

  const socketRef = useRef<WebSocket | null>(null);
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
    videosHandlers.setState([]);
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

      const newVideos = moveItem(videos, from, to).map((v, i) => ({ ...v, queuePosition: i + 1 }));

      socketRef.current?.send(JSON.stringify({ eventName: 'reorder', data: newVideos }));
    },
    [videos, socketRef.current],
  );

  useEffect(() => {
    if (!profile.data) return;

    if (!socketRef.current) {
      const url = `${`${window.location.protocol === 'https:' ? 'wss' : 'ws'}://${window.location.host}`}/socket/youtube?apiKey=${profile.data.apiKey}`;
      socketRef.current = new WebSocket(url);
    }

    socketRef.current!.onmessage = (msg) => {
      const event = JSON.parse(msg.data);

      if (event.eventName === 'currentQueue') {
        videosHandlers.setState([]);
        addVideos(event.data);
      }

      if (event.eventName === 'newTrack') {
        if (autoPlay === 0) setAutoPlay(1);
        addVideos([event.data]);
      }
    };

    return () => {
      socketRef.current?.close();
      socketRef.current = null;
    };
  }, [profile.data]);

  // it's in another useEffect because we need videos as dependency for correctly find index for skip
  useEffect(() => {
    if (!socketRef.current) return;

    socketRef.current?.addEventListener('message', (msg) => {
      const event = JSON.parse(msg.data);

      if (event.eventName === 'removeTrack') {
        const index = videos.findIndex((v) => v.id === event.data.id);
        if (index >= 0) {
          skipVideo(index, false);
        }
      }
    });
  }, [videos, socketRef.current]);

  function callWsSkip(videos: RequestedSong | RequestedSong[]) {
    const ids = (Array.isArray(videos) ? videos : [videos]).map((v) => v.id);
    socketRef.current?.send(JSON.stringify({
      eventName: 'skip',
      data: ids,
    }));
  }

  useEffect(() => {
    if (!socketRef.current) return;
    const video = videos[0]!;
    if (isPlaying) {
      socketRef.current?.send(JSON.stringify({
        eventName: 'play',
        data: { id: video.id, duration: video.duration },
      }));
    } else {
      socketRef.current?.send(JSON.stringify({
        eventName: 'pause',
      }));
    }
  }, [isPlaying]);

  return (
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
      <Grid>
        {!youtubeSettings?.enabled && <AlertPlayerDisabled />}
        {videos.length === 0 && <AlertQueueEmpty />}
        <Grid.Col md={4} lg={4}>
          <YoutubePlayer />
        </Grid.Col>
        <Grid.Col md={8} lg={8}>
          <QueueList />
        </Grid.Col>
      </Grid>
    </PlayerContext.Provider>
  );
};

export default Player;
