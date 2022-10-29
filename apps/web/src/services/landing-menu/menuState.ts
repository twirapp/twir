import { useStore } from '@nanostores/vue';
import { atom } from 'nanostores';

export const menuStateStore = atom<boolean>(false);

menuStateStore.listen((menuState) => {
  if (menuState) {
    document.body.style.overflow = 'hidden';
  } else {
    document.body.style.overflow = 'visible';
  }
});

export const useLandingMenuState = () => {
  const menuState = useStore(menuStateStore);

  const toggleMenuState = () => menuStateStore.set(!menuStateStore.get());

  const closeMenu = () => menuStateStore.set(false);

  return { menuState, toggleMenuState, closeMenu };
};
