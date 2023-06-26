import { cssPX, cssURL } from '@twir/ui-components';
import { useImage, useMouseInElement, type MaybeElementRef } from '@vueuse/core';
import { computed, ComputedRef, StyleValue } from 'vue';

export default function (
  cardRef: MaybeElementRef,
  imagePath: string,
  scaleImageIndex: number = 1,
): { styles: ComputedRef<StyleValue>; isActive: ComputedRef<boolean> } {
  const { elementX, elementY, isOutside } = useMouseInElement(cardRef);
  const { state, isReady } = useImage({
    src: imagePath,
  });

  const blobSize = computed<{ width: number; height: number }>(() => {
    if (isReady.value && state.value) {
      return {
        height: state.value.height * scaleImageIndex,
        width: state.value.width * scaleImageIndex,
      };
    }
    return { height: 0, width: 0 };
  });

  const blobPosition = computed<{ top: number; left: number }>(() => {
    if (!isOutside.value) {
      return {
        top: elementY.value - blobSize.value.height / 2,
        left: elementX.value - blobSize.value.width / 2,
      };
    }

    return { top: 0, left: 0 };
  });

  const bgImageUrl = cssURL(imagePath);
  const styles = computed(() => ({
    backgroundImage: bgImageUrl,
    height: cssPX(blobSize.value.height),
    width: cssPX(blobSize.value.width),
    top: cssPX(blobPosition.value.top),
    left: cssPX(blobPosition.value.left),
  }));

  const isActive = computed(() => !isOutside.value && isReady.value);

  return { styles, isActive };
}
