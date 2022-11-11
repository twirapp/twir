import { Socket } from 'socket.io-client';

import { usePlyrYoutubeQueue } from './usePlyrYoutubeQueue.js';

import { NAMESPACES, nameSpaces } from '@/plugins/socket.js';

export type RequestedSong = {
  id: string;
  videoId: string;
  title: string;
  duration: number;
  createdAt: Date;
  orderedById: string;
  orderedByName: string;
  orderedBy?: string;
  channelId: string;
  channel?: string;
  deletedAt: Date | null;
};

export const useYoutubeSocketPlayer = () => {
  const player = usePlyrYoutubeQueue([], {});

  interface SocketEvents {
    play: (video: { id: string }) => void;
    skip: (id: string) => void;
    newTrack: (video: RequestedSong) => void;
    currentQueue: (callback: (videos: RequestedSong[]) => void) => void;
  }

  const socket: Socket<SocketEvents> | undefined = nameSpaces.get(NAMESPACES.YOUTUBE);

  if (!socket) {
    throw new Error('Cannot get youtube socket');
  }

  player.onPlayVideo((video) => {
    console.log('play: ', video.id);
    socket.emit('play', { id: video.id });
  });
  player.onRemoveVideo((video) => {
    console.log('remove: ', video.id);
    socket.emit('skip', video.id);
  });

  socket.on('newTrack', (video: RequestedSong) => {
    console.log('newTrack: ', video.id);
    player.addVideo(video);
  });

  socket.emit('currentQueue', (videos: RequestedSong[]) => {
    if (videos.length === 0) return;
    player.addVideo(...videos);
  });

  return { ...player };
};
