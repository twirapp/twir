'use client';

import { Table } from '@mantine/core';
import type { RequestedSong } from '@tsuwari/typeorm/entities/RequestedSong';
import { getCookie } from 'cookies-next';
import { GetServerSideProps, NextPage } from 'next';
import { serverSideTranslations } from 'next-i18next/serverSideTranslations';
import dynamic from 'next/dynamic';
import { useCallback, useEffect, useRef, useState } from 'react';
import { io, Socket } from 'socket.io-client';

import { PlayerContext } from '@/components/song-requests/context';
import { useProfile } from '@/services/api';

export const getServerSideProps: GetServerSideProps = async ({ locale }) => ({
  props: {
    ...(await serverSideTranslations(locale!, ['layout'])),
  },
});

const PlayerComponent = dynamic(
  () => import('../../components/song-requests/player'),
  { ssr: false },
);

const Player: NextPage = () => {
  const profile = useProfile();
  const [videos, setVideos] = useState<RequestedSong[]>([]);
  const socketRef = useRef<Socket | null>(null);
  const [isPlaying, setIsPlaying] = useState(false);

  const skipVideo = useCallback(
    (index = 0) => {
      setIsPlaying(false);
      callWsSkip(videos[0]!);

      const length = videos.length;
      if (index === 0) {
        setVideos(videos.slice(1));
      } else if (index === length - 1) {
        setVideos(videos.slice(0, length - 1));
      } else {
        setVideos([...videos.slice(0, index), ...videos.slice(index + 1)]);
      }
    },
    [videos],
  );

  const addVideos = useCallback(
    (v: RequestedSong[]) => {
      setVideos([...videos, ...v]);
    },
    [videos],
  );

  useEffect(() => {
    if (!profile.data) return;

    if (!socketRef.current) {
      socketRef.current = io(`${`${window.location.protocol == 'https:' ? 'wss' : 'ws'}://${window.location.host}`}/youtube`, {
        transports: ['websocket'],
        autoConnect: false,
        auth: (cb) => {
          cb({ apiKey: profile.data?.apiKey, channelId: getCookie('selectedDashboard') });
        },
      });
    }

    socketRef.current.connect();

    socketRef.current.emit('currentQueue', (data: RequestedSong[]) => {
      addVideos(data);
    });

    socketRef.current.on('newTrack', (track: RequestedSong) => {
      addVideos([track]);
    });

    socketRef.current.on('removeTrack', (track: RequestedSong) => {
      const index = videos.findIndex(v => v.id === track.id);
      if (index) {
        skipVideo(index);
      }
    });

    return () => {
      socketRef.current?.off('newTrack');
      socketRef.current?.off('removeTrack');
      socketRef.current?.disconnect();
    };
  }, [profile.data]);

  function callWsSkip(video: RequestedSong) {
    socketRef.current?.emit('skip', video.id);
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
    <div>
      <PlayerContext.Provider
        value={{
          videos,
          skipVideo,
          addVideos,
          isPlaying,
          setIsPlaying,
        }}
      ><PlayerComponent/></PlayerContext.Provider>

      <Table>
        <table>
          <thead>
          <tr>
            <th>#</th>
            <th>Title</th>
            <th>Requested by</th>
          </tr>
          </thead>
          <tbody>
          {videos?.map((video, index) => (
            <tr key={video.id}>
              <th>{index + 1}</th>
              <th><a href={'https://youtu.be/' + video.videoId}>{video.title}</a></th>
              <th>{video.orderedByName}</th>
            </tr>
          ))}
          </tbody>
        </table>
      </Table>
    </div>
  );
};

export default Player;