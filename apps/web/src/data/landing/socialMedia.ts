export type SocialMedia = 'Twitch' | 'Telegram' | 'Github' | 'Instagram' | 'Website';

export interface SocialMediaItem {
  id: number;
  type: SocialMedia;
  href: string;
}

export const socials: SocialMediaItem[] = [
  { id: 1, type: 'Telegram', href: '#' },
  { id: 2, type: 'Instagram', href: '#' },
];
