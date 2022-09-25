import { atom } from 'nanostores';

type Profile = {
  id: number;
  display_name: string;
  username: string;
  thumbnail: string;
};

export const streamlabsStore = atom<Profile | null | undefined>(null);

export function setStreamlabsStore(data: Profile | null) {
  streamlabsStore.set(data);
}
