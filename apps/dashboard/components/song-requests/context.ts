import { UseListStateHandlers } from '@mantine/hooks';
import type { RequestedSong } from '@tsuwari/typeorm/entities/RequestedSong';
import { createContext, Dispatch, SetStateAction } from 'react';
import { DraggableLocation } from 'react-beautiful-dnd';

export const PlayerContext = createContext({
  videos: [] as RequestedSong[],
} as {
  videos: RequestedSong[];
  videosHandlers: UseListStateHandlers<RequestedSong>;
  reorderVideos: (destination: DraggableLocation, source: DraggableLocation) => void;
  addVideos: (v: RequestedSong[]) => void;
  skipVideo: (index?: number) => void;
  isPlaying: boolean;
  setIsPlaying: Dispatch<SetStateAction<boolean>>;
  autoPlay: number;
  setAutoPlay: Dispatch<SetStateAction<number>>;
});
