import { type Profile } from '@twir/grpc/generated/api/api/auth';
import { atom } from 'nanostores';

export const profileStore = atom<Profile | null>(null);
export const authLinkStore = atom('');
