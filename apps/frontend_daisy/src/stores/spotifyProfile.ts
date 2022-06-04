import { atom } from 'nanostores';

type Profile = {
  id: string;
  display_name: string;
};

export const spotifyProfileStore = atom<Profile | null | undefined>(null);

export function setSpotifyProfile(data: Profile | null) {
  spotifyProfileStore.set(data);
}

