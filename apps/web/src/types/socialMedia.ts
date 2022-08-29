export type SocialMedia = 'Twitch' | 'Telegram' | 'Github' | 'Instagram' | 'Website';

export interface SocialMediaItem {
  id: number;
  media: SocialMedia;
  href: string;
}
