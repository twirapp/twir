import type { RequestedSong } from '@tsuwari/typeorm/entities/RequestedSong';
import { createContext, Dispatch, SetStateAction } from 'react';

export const PlayerContext = createContext({
  videos: [],
  setVideos: () => {
  },
} as {
  videos: RequestedSong[],
  setVideos: Dispatch<SetStateAction<RequestedSong[]>>,
});