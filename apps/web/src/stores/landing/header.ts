import { atom, computed } from 'nanostores';

export const headerStore = atom<HTMLElement | null>(null);

export const headerHeightStore = computed([headerStore], (h) => (h !== null ? h.clientHeight : 0));

export const menuStateStore = atom<boolean>(false);
