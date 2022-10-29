import { useImage, useIntersectionObserver } from '@vueuse/core';
import { computed, Ref, ref } from 'vue';

export default function (
  imageSrc: string,
  placeholder: Ref<HTMLElement | undefined>,
  immediate: boolean = true,
) {
  const isTargetVisible = ref<boolean>(false);
  const {
    isReady,
    execute: executeImage,
    state: image,
  } = useImage({ src: imageSrc }, { immediate: false });

  const execute = () => {
    const { stop } = useIntersectionObserver(placeholder, async ([{ isIntersecting }]) => {
      if (isIntersecting) {
        isTargetVisible.value = true;
        stop();
        executeImage();
      }
    });
  };

  const isImageReady = computed(() => isTargetVisible.value && isReady.value);

  if (immediate) {
    execute();
    return { isImageReady, image };
  }

  return { isImageReady, image, execute };
}
