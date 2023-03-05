import type { SocialMedia } from '@/data/landing/socialMedia';

export enum TeamMemberId {
  'Satont',
  'Melkam',
}

export type TeamMemberMedia = {
  type: SocialMedia;
  link: string;
};

export interface TeamMember {
  name: string;
  isFounder?: boolean;
  socials: TeamMemberMedia[];
  avatarUrl?: string,
}

type TeamMembers = {
  [K in TeamMemberId]: TeamMember;
};

export type TeamMemberLocale = {
  [K in TeamMemberId]: string; // string is member role
};

export const teamMembers: TeamMembers = {
  [TeamMemberId.Satont]: {
    name: 'Satont',
    isFounder: true,
    socials: [
      { type: 'Twitch', link: 'https://www.twitch.tv/fukushine' },
      { type: 'Telegram', link: 'https://t.me/satont' },
      { type: 'Github', link: 'https://github.com/satont' },
      { type: 'Website', link: 'https://satont.dev/' },
    ],
    avatarUrl: 'https://cdn.7tv.app/emote/62c5c34724fb1819d9f08b4d/3x.webp',
  },
  [TeamMemberId.Melkam]: {
    name: 'Melkam',
    socials: [
      { type: 'Twitch', link: 'https://www.twitch.tv/mellkam' },
      { type: 'Telegram', link: 'https://t.me/mellkam' },
      { type: 'Github', link: 'https://github.com/MellKam' },
      { type: 'Instagram', link: 'https://www.instagram.com/mel._.kam/' },
    ],
    avatarUrl: 'https://avatars.githubusercontent.com/u/51422045?s=80&v=4',
  },
};
