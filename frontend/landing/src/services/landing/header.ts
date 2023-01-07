import { useStore } from '@nanostores/vue';
import { atom, computed } from 'nanostores';
import type { ComponentPublicInstance } from 'vue';

export const headerStore = atom<HTMLElement | null>(null);

export const headerHeightStore = computed([headerStore], (h) => (h !== null ? h.clientHeight : 0));

export const useLandingHeaderHeight = () => useStore(headerHeightStore);

export const setHeaderRef = (h: Element | ComponentPublicInstance | null) => {
  if (h === null) return;
  headerStore.set(h as HTMLElement);
};
