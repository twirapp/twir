import { Group, Slider, Text } from '@mantine/core';
import { useInterval } from '@mantine/hooks';
import { useState, useCallback, useEffect } from 'react';
import { YouTubePlayer } from 'react-youtube';

import { convertMillisToTime } from '../helpers';

export function PlayerSlider({
  isPlaying,
  player,
}: {
  isPlaying: boolean;
  player: YouTubePlayer | null;
}) {
  const [currentTime, setCurrentTime] = useState(0);
  const currentTimeInterval = useInterval(() => getSongCurrentTime(), 1000);
  const getSongCurrentTime = useCallback(() => {
    if (!player) return setCurrentTime(0);
    setCurrentTime(player?.getCurrentTime() as unknown as number);
  }, [player]);

  const [songDuration, setSongDuration] = useState(0);
  const songDurationInterval = useInterval(() => getSongDuration(), 1000);
  const getSongDuration = useCallback(() => {
    if (!player) return setSongDuration(0);
    setSongDuration(player?.getDuration() as unknown as number);
  }, [player]);

  const setTime = useCallback(
    (time: number) => {
      player?.seekTo(time, true);
    },
    [currentTime, player],
  );

  useEffect(() => {
    currentTimeInterval.start();
    songDurationInterval.start();

    getSongDuration();
    getSongCurrentTime();

    return () => {
      currentTimeInterval.stop();
      songDurationInterval.stop();
    };
  }, [isPlaying]);

  return (
    <>
      <Group position="apart">
        <Text size="sm">{convertMillisToTime(currentTime * 1000)}</Text>
        <Text size="sm">{convertMillisToTime(songDuration * 1000)}</Text>
      </Group>
      <Slider
        step={1}
        labelTransition="pop"
        labelTransitionDuration={200}
        labelTransitionTimingFunction="ease"
        value={parseInt(currentTime.toFixed(0), 10)}
        label={(v) => convertMillisToTime(v * 1000)}
        onChange={(v) => {
          if (currentTimeInterval.active) {
            currentTimeInterval.stop();
          }

          setCurrentTime(v);
        }}
        onChangeEnd={(v) => {
          currentTimeInterval.start();
          setTime(v);
        }}
        max={songDuration}
      />
    </>
  );
}
