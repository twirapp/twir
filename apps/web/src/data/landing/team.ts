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
      { type: 'Github', link: 'https://github.com/satont' },
      { type: 'Telegram', link: 'https://t.me/satont' },
      { type: 'Twitch', link: 'https://www.twitch.tv/sadisnamenya' },
      { type: 'Website', link: 'https://satont.dev/' },
    ],
  },
  [TeamMemberId.Melkam]: {
    name: 'Melkam',
    socials: [
      { type: 'Github', link: 'https://github.com/MellKam' },
      { type: 'Instagram', link: 'https://www.instagram.com/mel._.kam/' },
      { type: 'Telegram', link: 'https://t.me/mellkam' },
    ],
  },
};
