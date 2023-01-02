import { UseListStateHandlers } from '@mantine/hooks';
import type { RequestedSong } from '@tsuwari/typeorm/entities/RequestedSong';
import { createContext, Dispatch, SetStateAction } from 'react';

export const PlayerContext = createContext({
  videos: [] as RequestedSong[],
} as {
  videos: RequestedSong[],
  videosHandlers:  UseListStateHandlers<RequestedSong>,
  addVideos: (v: RequestedSong[]) => void,
  skipVideo: (index?: number) => void,
  isPlaying: boolean,
  setIsPlaying: Dispatch<SetStateAction<boolean>>
});