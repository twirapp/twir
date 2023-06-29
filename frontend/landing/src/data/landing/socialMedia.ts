import type { IconName } from '@twir/ui-components';

type Extract<T, U extends T> = T extends U ? T : never;

export type SocialMedia = Extract<
  IconName,
  'Twitch' | 'Instagram' | 'Website' | 'Github' | 'Telegram' | 'Discord'
>;

export interface SocialMediaItem {
  id: number;
  type: SocialMedia;
  href: string;
}

export const socials: SocialMediaItem[] = [
  { id: 1, type: 'Telegram', href: 'https://t.me/tsuwari_app' },
  { id: 2, type: 'Discord', href: 'https://discord.gg/Q9NBZq3zVV' },
];
