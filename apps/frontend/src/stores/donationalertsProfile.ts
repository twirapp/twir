import { atom } from 'nanostores';

type Profile = {
  name: string;
  code: string;
  avatar: string;
};

export const donationAlertsStore = atom<Profile | null | undefined>(null);

export function setDonationAlertsStore(data: Profile | null) {
  donationAlertsStore.set(data);
}
