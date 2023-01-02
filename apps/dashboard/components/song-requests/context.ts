import type { RequestedSong } from '@tsuwari/typeorm/entities/RequestedSong';
import { createContext, Dispatch, SetStateAction } from 'react';

export const PlayerContext = createContext({
  videos: [] as RequestedSong[],
} as {
  videos: RequestedSong[],
  addVideos: (v: RequestedSong[]) => void,
  skipVideo: (index?: number) => void,
  isPlaying: boolean,
  setIsPlaying: Dispatch<SetStateAction<boolean>>
  setVideos: Dispatch<SetStateAction<RequestedSong[]>>
});