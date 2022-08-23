export type TeamMemberMedia = 'Twitch' | 'Telegram' | 'Github' | 'Instagram' | 'Website';

export type TeamMemberMediaLink = {
  media: TeamMemberMedia;
  link: string;
};

export interface TeamMember {
  id: number;
  name: string;
  role: string;
  isFounder?: boolean;
  socials: TeamMemberMediaLink[];
}
