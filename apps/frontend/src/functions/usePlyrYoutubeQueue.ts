import { createEventHook } from '@vueuse/core';
import type Plyr from 'plyr';
import { computed, onBeforeUnmount, readonly, ref, watch } from 'vue';

import { RequestedSong as Video } from './useYoutubeSocketPlayer.js';

export const usePlyrYoutubeQueue = (
  initialQueue: Video[],
  options: { autoplay?: boolean } = { autoplay: true },
) => {
  const player = ref<Plyr | null>(null);
  const autoplay = ref<boolean>(typeof options.autoplay === 'boolean' ? options.autoplay : true);
  const queue = ref<Video[]>(initialQueue);
  const isPaused = ref(true);
  const queueLength = computed(() => queue.value.length);
  const isActive = computed<boolean>(() => queueLength.value > 0);
  const playVideoEvent = createEventHook<{ video: Video; timeToEnd: number }>();
  const removeVideoEvent = createEventHook<Video>();
  const videoPauseEvent = createEventHook<Video>();
  const currentVideo = computed<Video | undefined>(() => queue.value[0]);
  const queueWithoutFirst = computed(() => {
    if (queueLength.value <= 1) {
      return [];
    }

    return [...queue.value].slice(1);
  });

  const addVideo = (...videos: Video[]) => {
    if (queue.value.length === 0) {
      setVideo(videos[0], false);
    }
    queue.value.push(...videos);
  };

  const removeVideo = (id: string) => {
    if (queue.value[0].id === id) {
      return;
    }
    queue.value = queue.value.filter((v) => {
      if (v.id !== id) {
        return true;
      }
      removeVideoEvent.trigger(v);
      return false;
    });
  };

  const initQueue = (plyr: Plyr) => {
    player.value = plyr;
    player.value.on('ended', next);

    if (queueLength.value > 0) {
      setVideo(queue.value[0], false);
    }

    const unwatchIsPaused = watch(
      () => isPaused.value,
      (isPaused) => {
        if (!isActive.value) return;
        if (player.value === null) {
          return console.error('Cannot pause video, because player in null');
        }
        if (isPaused && !player.value.paused) {
          return player.value.pause();
        }
        if (!isPaused && player.value.paused) {
          player.value.play();
        }
      },
    );
    onBeforeUnmount(unwatchIsPaused);

    player.value.on('pause', () => {
      if (currentVideo.value === undefined) return;

      isPaused.value = true;
      videoPauseEvent.trigger(currentVideo.value);
    });
    player.value.on('play', () => {
      if (!player.value || !currentVideo.value) return;

      isPaused.value = false;
      playVideoEvent.trigger({
        video: currentVideo.value,
        timeToEnd: Math.floor(currentVideo.value.duration - player.value.currentTime * 1000),
      });
    });
  };

  const setVideo = (video: Video, playImmediately?: boolean) => {
    if (player.value === null) {
      return console.error('Cannot set video, because player is null');
    }
    player.value.source = {
      type: 'video',
      sources: [
        {
          src: video.videoId,
          provider: 'youtube',
        },
      ],
    };

    playImmediately = typeof playImmediately === 'undefined' ? autoplay.value : playImmediately;

    if (playImmediately) {
      player.value.once('ready', () => {
        if (player.value === null) {
          return console.error('Cannot play video, because player is null');
        }
        player.value.currentTime = 0;
        player.value.play();
      });
    }
  };

  const skipCurrentVideo = () => {
    if (!isActive.value) return;

    removeVideoEvent.trigger(queue.value.shift()!);
    const nextVideo = queue.value[0];
    if (nextVideo) return setVideo(nextVideo);

    player.value?.stop();
  };

  function next() {
    if (queue.value.length === 0) {
      return player.value?.stop();
    }

    if (queue.value.length === 1) {
      return removeVideoEvent.trigger(queue.value.pop()!);
    }

    removeVideoEvent.trigger(queue.value.shift()!);
    setVideo(queue.value[0]);
  }

  return {
    addVideo,
    removeVideo,
    initQueue,
    isPaused,
    skipCurrentVideo,
    autoplay,
    isActive: readonly(isActive),
    queueLength: readonly(queueLength),
    queue: readonly(queue),
    onPlayVideo: playVideoEvent.on,
    onRemoveVideo: removeVideoEvent.on,
    onPause: videoPauseEvent.on,
    currentVideo: readonly(currentVideo),
    queueWithoutFirst: readonly(queueWithoutFirst),
  };
};
