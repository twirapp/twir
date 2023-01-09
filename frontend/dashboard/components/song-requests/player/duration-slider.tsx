import { Group, Slider, Text } from '@mantine/core';
import { useInterval } from '@mantine/hooks';
import { useState, useCallback, useEffect } from 'react';
import { YouTubePlayer } from 'react-youtube';

import { convertMillisToTime, toFixedNum } from '../helpers';

export function DurationSlider({
  isPlaying,
  player,
}: {
  isPlaying: boolean;
  player: YouTubePlayer | null;
}) {
  const [currentTime, setCurrentTime] = useState(0);
  const [songDuration, setSongDuration] = useState(0);

  const interval = useInterval(() => {
    getSongCurrentTime();
    getSongDuration();
  }, 1000);

  const getSongCurrentTime = useCallback(() => {
    if (!player) return setCurrentTime(0);
    setCurrentTime(player?.getCurrentTime());
  }, [player]);

  const getSongDuration = useCallback(() => {
    if (!player) return setSongDuration(0);
    setSongDuration(player?.getDuration());
  }, [player]);

  const updatePlayerTime = useCallback(
    (time: number) => {
      player?.seekTo(time, true);
    },
    [player],
  );

  useEffect(() => {
    interval.start();

    getSongDuration();
    getSongCurrentTime();

    return () => {
      interval.stop();
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
        size="sm"
        labelTransition="pop"
        labelTransitionDuration={200}
        labelTransitionTimingFunction="ease"
        value={toFixedNum(currentTime)}
        label={(v) => convertMillisToTime(v * 1000)}
        onChange={(v) => {
          if (interval.active) {
            interval.stop();
          }

          setCurrentTime(v);
        }}
        onChangeEnd={(v) => {
          interval.start();
          updatePlayerTime(v);
        }}
        max={songDuration}
      />
    </>
  );
}
