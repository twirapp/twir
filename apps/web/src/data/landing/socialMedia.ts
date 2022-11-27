import type { IconName } from '@tsuwari/ui-components';

type Extract<T, U extends T> = T extends U ? T : never;

export type SocialMedia = Extract<
  IconName,
  'Twitch' | 'Instagram' | 'Website' | 'Github' | 'Telegram'
>;

export interface SocialMediaItem {
  id: number;
  type: SocialMedia;
  href: string;
}

export const socials: SocialMediaItem[] = [
  { id: 1, type: 'Telegram', href: '#' },
  { id: 2, type: 'Instagram', href: '#' },
];
