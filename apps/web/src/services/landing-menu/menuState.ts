import { useStore } from '@nanostores/vue';
import { isClient } from '@vueuse/core';
import { atom } from 'nanostores';
import { onUnmounted } from 'vue';

export const menuStateStore = atom<boolean>(false);

export const useLandingMenuState = () => {
  const menuState = useStore(menuStateStore);

  const toggleMenuState = () => menuStateStore.set(!menuStateStore.get());

  const closeMenu = () => menuStateStore.set(false);

  if (menuStateStore.lc === 0 && isClient) {
    const removeListener = menuStateStore.listen((menuState) => {
      if (menuState) {
        document.body.style.overflow = 'hidden';
      } else {
        document.body.style.overflow = 'visible';
      }
    });

    onUnmounted(() => {
      removeListener();
    });
  }

  return { menuState, toggleMenuState, closeMenu };
};
